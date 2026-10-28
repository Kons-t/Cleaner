package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/huh/v2/spinner"

	"cleaner/internal/domain/enums"
	"cleaner/internal/domain/paths"
	"cleaner/internal/usecase/davinci_clean"
	"cleaner/internal/usecase/docker_clean"
	"cleaner/internal/usecase/ds_store_clean"
	"cleaner/internal/usecase/go_clean"
	"cleaner/internal/usecase/homebrew_clean"
	"cleaner/internal/usecase/idea_http_clean"
	"cleaner/internal/usecase/logs_clean"
	"cleaner/internal/usecase/xcode_clean"
)

const (
	exitOption = "Выход"
	exitChoice = -1

	exitHint   = "Esc / Ctrl+C — выход из программы"
	cancelHint = "Esc / Ctrl+C — отмена и возврат в главное меню"
)

// menuItem — пункт корневого меню: короткая метка для вступительного сообщения,
// заголовок пункта, обработчик и проверка, доступен ли модуль на этом устройстве.
// Новые модули очистки добавляются сюда, не трогая Run.
type menuItem struct {
	label     string
	title     string
	handler   func(ctx context.Context) error
	available func(ctx context.Context) bool
}

// Application — оркестрирует меню, ввод пользователя и вызовы usecase-слоя.
type Application struct {
	goCleaner            go_clean.Cleaner
	httpResponsesCleaner idea_http_clean.Cleaner
	homebrewCleaner      homebrew_clean.Cleaner
	dockerCleaner        docker_clean.Cleaner
	logsCleaner          logs_clean.Cleaner
	dsStoreCleaner       ds_store_clean.Cleaner
	davinciCleaner       davinci_clean.Cleaner
	xcodeCleaner         xcode_clean.Cleaner

	keyMap *huh.KeyMap
}

// New создаёт Application с переданными зависимостями usecase-слоя.
func New(
	goCleaner go_clean.Cleaner,
	httpResponsesCleaner idea_http_clean.Cleaner,
	homebrewCleaner homebrew_clean.Cleaner,
	dockerCleaner docker_clean.Cleaner,
	logsCleaner logs_clean.Cleaner,
	dsStoreCleaner ds_store_clean.Cleaner,
	davinciCleaner davinci_clean.Cleaner,
	xcodeCleaner xcode_clean.Cleaner,
) *Application {
	return &Application{
		goCleaner:            goCleaner,
		httpResponsesCleaner: httpResponsesCleaner,
		homebrewCleaner:      homebrewCleaner,
		dockerCleaner:        dockerCleaner,
		logsCleaner:          logsCleaner,
		dsStoreCleaner:       dsStoreCleaner,
		davinciCleaner:       davinciCleaner,
		xcodeCleaner:         xcodeCleaner,
		keyMap:               newKeyMap(),
	}
}

// Run запускает главный цикл меню до тех пор, пока пользователь не выберет выход.
func (a *Application) Run(ctx context.Context) error {
	candidates := []menuItem{
		{label: "Go-кеши", title: "Go: очистить кеши", handler: a.runGoCacheMenu, available: a.goCleaner.Available},
		{
			label:     "HTTP-ответы IDEA",
			title:     fmt.Sprintf("IDEA: очистить .json ответы HTTP-клиента (%s)", paths.WorkDir),
			handler:   a.runHTTPResponsesMenu,
			available: a.httpResponsesCleaner.Available,
		},
		{
			label:     "Homebrew-кеш",
			title:     "Homebrew: очистить кеш (brew cleanup)",
			handler:   a.runHomebrewMenu,
			available: a.homebrewCleaner.Available,
		},
		{
			label:     "Docker-ресурсы",
			title:     "Docker: очистить неиспользуемые ресурсы (docker system prune)",
			handler:   a.runDockerMenu,
			available: a.dockerCleaner.Available,
		},
		{
			label:     "логи приложений",
			title:     fmt.Sprintf("Логи: очистить файлы приложений (%s)", paths.LogsDir),
			handler:   a.runLogsMenu,
			available: a.logsCleaner.Available,
		},
		{
			label:     ".DS_Store",
			title:     fmt.Sprintf(".DS_Store: удалить мусорные файлы Finder (%s)", paths.DSStoreRoot),
			handler:   a.runDSStoreMenu,
			available: a.dsStoreCleaner.Available,
		},
		{
			label:     "render cache DaVinci Resolve",
			title:     fmt.Sprintf("DaVinci Resolve: очистить render cache (%s)", paths.DavinciCacheDir),
			handler:   a.runDavinciMenu,
			available: a.davinciCleaner.Available,
		},
		{
			label:     "Xcode DerivedData",
			title:     fmt.Sprintf("Xcode: очистить DerivedData (%s)", paths.XcodeDerivedDataDir),
			handler:   a.runXcodeMenu,
			available: a.xcodeCleaner.Available,
		},
	}

	items := make([]menuItem, 0, len(candidates))
	labels := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if !c.available(ctx) {
			continue
		}
		items = append(items, c)
		labels = append(labels, c.label)
	}

	if len(labels) == 0 {
		fmt.Println("\nНа этом устройстве не найдено ни одного поддерживаемого инструмента для очистки.")
	} else {
		fmt.Printf(
			"\nЭТО ПРИЛОЖЕНИЕ ДЛЯ УДАЛЕНИЯ МУСОРНЫХ ФАЙЛОВ\n\nОно умеет удалять:\n%s\n— показаны только инструменты, найденные на этом устройстве.\n\n",
			strings.Join(labels, ", "),
		)
	}

	options := make([]huh.Option[int], 0, len(items)+1)
	for i, item := range items {
		options = append(options, huh.NewOption(item.title, i))
	}
	options = append(options, huh.NewOption(exitOption, exitChoice))

	for {
		choice := exitChoice
		err := a.runField(ctx, huh.NewSelect[int]().
			Title("Выберите действие:").
			Description(exitHint).
			Options(options...).
			Value(&choice),
		)
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return err
		}

		if choice == exitChoice {
			return nil
		}

		if err := items[choice].handler(ctx); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				fmt.Println("Отменено.")
				continue
			}
			fmt.Println("Ошибка:", err)
		}
	}
}

func (a *Application) runGoCacheMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Считаю размеры кешей, это может занять несколько секунд...", a.goCleaner.Status)
	if err != nil {
		return err
	}

	fmt.Printf("GOCACHE:    %s (%s)\n", status.GoCache.Path, status.GoCache.SizeHuman)
	fmt.Printf("GOMODCACHE: %s (%s)\n", status.GoModCache.Path, status.GoModCache.SizeHuman)

	types := enums.AllCacheTypes()
	options := make([]huh.Option[enums.CacheType], 0, len(types))
	for _, t := range types {
		options = append(options, huh.NewOption(t.Label(), t))
	}

	var selected []enums.CacheType
	if err := a.runField(ctx, huh.NewMultiSelect[enums.CacheType]().
		Title("Что очистить (Space — отметить, Ctrl+A — выбрать всё, Enter — подтвердить):").
		Description(cancelHint).
		Options(options...).
		Value(&selected),
	); err != nil {
		return err
	}

	if len(selected) == 0 {
		return nil
	}

	confirmed, err := a.confirm(ctx, "Точно очистить выбранное? Действие необратимо.")
	if err != nil || !confirmed {
		return err
	}

	results := a.goCleaner.Clean(ctx, selected)
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("✗ %s: ошибка — %v\n", r.Type.Label(), r.Err)
			continue
		}
		fmt.Printf("✓ %s очищен\n", r.Type.Label())
	}

	return nil
}

func (a *Application) runHTTPResponsesMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Ищу .json-ответы HTTP-клиента в .idea/httpRequests...", a.httpResponsesCleaner.Status)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed, err := a.confirm(ctx, fmt.Sprintf("Удалить все %d .json-файлов? Действие необратимо.", status.FileCount))
	if err != nil || !confirmed {
		return err
	}

	result := a.httpResponsesCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runHomebrewMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Смотрю, что можно почистить (brew cleanup --dry-run)...", a.homebrewCleaner.Status)
	if err != nil {
		return err
	}

	fmt.Println(status.Preview)

	confirmed, err := a.confirm(ctx, "Выполнить brew cleanup? Действие необратимо.")
	if err != nil || !confirmed {
		return err
	}

	result, err := a.homebrewCleaner.Clean(ctx)
	if err != nil {
		return err
	}

	fmt.Println(result.Output)

	return nil
}

func (a *Application) runDockerMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Смотрю использование диска Docker (docker system df)...", a.dockerCleaner.Status)
	if err != nil {
		return err
	}

	fmt.Println(status.DiskUsage)

	confirmed, err := a.confirm(
		ctx,
		"Выполнить docker system prune -f (остановленные контейнеры, висячие образы, неиспользуемые сети, build cache)? Действие необратимо.",
	)
	if err != nil || !confirmed {
		return err
	}

	result, err := a.dockerCleaner.Clean(ctx)
	if err != nil {
		return err
	}

	fmt.Println(result.Output)

	return nil
}

func (a *Application) runLogsMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Ищу файлы логов...", a.logsCleaner.Status)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed, err := a.confirm(ctx, fmt.Sprintf("Удалить все %d файлов логов? Действие необратимо.", status.FileCount))
	if err != nil || !confirmed {
		return err
	}

	result := a.logsCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runDSStoreMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Ищу .DS_Store файлы...", a.dsStoreCleaner.Status)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed, err := a.confirm(ctx, fmt.Sprintf("Удалить все %d файлов .DS_Store? Действие необратимо.", status.FileCount))
	if err != nil || !confirmed {
		return err
	}

	result := a.dsStoreCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runDavinciMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Ищу файлы render cache DaVinci Resolve...", a.davinciCleaner.Status)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed, err := a.confirm(ctx, fmt.Sprintf("Удалить все %d файлов render cache? Действие необратимо.", status.FileCount))
	if err != nil || !confirmed {
		return err
	}

	result := a.davinciCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runXcodeMenu(ctx context.Context) error {
	status, err := withSpinner(ctx, "Ищу папки проектов в DerivedData...", a.xcodeCleaner.Status)
	if err != nil {
		return err
	}

	if status.EntryCount == 0 {
		fmt.Println("Папок не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено папок: %d (%s) в %s\n", status.EntryCount, status.SizeHuman, status.Root)

	confirmed, err := a.confirm(ctx, fmt.Sprintf("Удалить все %d папок DerivedData? Действие необратимо.", status.EntryCount))
	if err != nil || !confirmed {
		return err
	}

	result := a.xcodeCleaner.Clean(ctx)
	fmt.Printf("Удалено папок: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

// confirm спрашивает подтверждение с кнопками «Да»/«Нет».
// При отмене через Esc / Ctrl+C возвращает huh.ErrUserAborted.
func (a *Application) confirm(ctx context.Context, message string) (bool, error) {
	confirmed := false
	err := a.runField(ctx, huh.NewConfirm().
		Title(message).
		Description(cancelHint).
		Affirmative("Да").
		Negative("Нет").
		Value(&confirmed),
	)

	return confirmed, err
}

// runField запускает одиночное поле как форму с общим keymap приложения.
func (a *Application) runField(ctx context.Context, field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(a.keyMap).
		RunWithContext(ctx)
}

// newKeyMap — стандартный keymap huh, в котором выход/отмена доступны и по Esc.
// Фильтр по «/» отключён: в нём Esc сбрасывает фильтр и конфликтовал бы с выходом,
// а в коротких списках меню фильтр не нужен.
func newKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("esc", "ctrl+c"))
	km.Select.Filter.SetEnabled(false)
	km.MultiSelect.Filter.SetEnabled(false)

	return km
}

// withSpinner выполняет action, показывая спиннер с заголовком title.
// Прерывание спиннера по Ctrl+C возвращается как huh.ErrUserAborted.
// action при этом не останавливается сразу, поэтому спиннер используется
// только для операций чтения (Status), но не для удаления.
func withSpinner[T any](ctx context.Context, title string, action func(ctx context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var result T
	err := spinner.New().
		Title(title).
		Context(ctx).
		ActionWithErr(func(ctx context.Context) error {
			var actionErr error
			result, actionErr = action(ctx)

			return actionErr
		}).
		Run()
	if err != nil {
		var zero T
		if errors.Is(err, tea.ErrInterrupted) {
			return zero, huh.ErrUserAborted
		}

		return zero, err
	}

	return result, nil
}
