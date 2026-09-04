package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

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

const exitOption = "Выход"

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

	options := make([]string, 0, len(items)+1)
	for _, item := range items {
		options = append(options, item.title)
	}
	options = append(options, exitOption)

	for {
		var choice string
		prompt := &survey.Select{
			Message: "Выберите действие:",
			Options: options,
		}
		if err := survey.AskOne(prompt, &choice); err != nil {
			return err
		}

		if choice == exitOption {
			return nil
		}

		for _, item := range items {
			if item.title == choice {
				if err := item.handler(ctx); err != nil {
					fmt.Println("Ошибка:", err)
				}
				break
			}
		}
	}
}

func (a *Application) runGoCacheMenu(ctx context.Context) error {
	fmt.Println("Считаю размеры кешей, это может занять несколько секунд...")

	status, err := a.goCleaner.Status(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("GOCACHE:    %s (%s)\n", status.GoCache.Path, status.GoCache.SizeHuman)
	fmt.Printf("GOMODCACHE: %s (%s)\n", status.GoModCache.Path, status.GoModCache.SizeHuman)

	const allOption = "Всё сразу"

	types := enums.AllCacheTypes()
	options := make([]string, 0, len(types)+1)
	for _, t := range types {
		options = append(options, t.Label())
	}
	options = append(options, allOption)

	var chosen []string
	prompt := &survey.MultiSelect{
		Message: "Что очистить (Space — отметить, Enter — подтвердить):",
		Options: options,
	}
	if err := survey.AskOne(prompt, &chosen); err != nil {
		return err
	}

	if len(chosen) == 0 {
		return nil
	}

	selected := resolveSelectedTypes(chosen, types, allOption)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Точно очистить выбранное? Действие необратимо.",
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
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
	fmt.Println("Ищу .json-ответы HTTP-клиента в .idea/httpRequests...")

	status, err := a.httpResponsesCleaner.Status(ctx)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Удалить все %d .json-файлов? Действие необратимо.", status.FileCount),
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result := a.httpResponsesCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runHomebrewMenu(ctx context.Context) error {
	fmt.Println("Смотрю, что можно почистить (brew cleanup --dry-run)...")

	status, err := a.homebrewCleaner.Status(ctx)
	if err != nil {
		return err
	}

	fmt.Println(status.Preview)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Выполнить brew cleanup? Действие необратимо.",
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result, err := a.homebrewCleaner.Clean(ctx)
	if err != nil {
		return err
	}

	fmt.Println(result.Output)

	return nil
}

func (a *Application) runDockerMenu(ctx context.Context) error {
	fmt.Println("Смотрю использование диска Docker (docker system df)...")

	status, err := a.dockerCleaner.Status(ctx)
	if err != nil {
		return err
	}

	fmt.Println(status.DiskUsage)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Выполнить docker system prune -f (остановленные контейнеры, висячие образы, неиспользуемые сети, build cache)? Действие необратимо.",
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result, err := a.dockerCleaner.Clean(ctx)
	if err != nil {
		return err
	}

	fmt.Println(result.Output)

	return nil
}

func (a *Application) runLogsMenu(ctx context.Context) error {
	fmt.Println("Ищу файлы логов...")

	status, err := a.logsCleaner.Status(ctx)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Удалить все %d файлов логов? Действие необратимо.", status.FileCount),
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result := a.logsCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runDSStoreMenu(ctx context.Context) error {
	fmt.Println("Ищу .DS_Store файлы...")

	status, err := a.dsStoreCleaner.Status(ctx)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Удалить все %d файлов .DS_Store? Действие необратимо.", status.FileCount),
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result := a.dsStoreCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runDavinciMenu(ctx context.Context) error {
	fmt.Println("Ищу файлы render cache DaVinci Resolve...")

	status, err := a.davinciCleaner.Status(ctx)
	if err != nil {
		return err
	}

	if status.FileCount == 0 {
		fmt.Println("Файлов не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено файлов: %d (%s) в %s\n", status.FileCount, status.SizeHuman, status.Root)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Удалить все %d файлов render cache? Действие необратимо.", status.FileCount),
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result := a.davinciCleaner.Clean(ctx)
	fmt.Printf("Удалено файлов: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func (a *Application) runXcodeMenu(ctx context.Context) error {
	fmt.Println("Ищу папки проектов в DerivedData...")

	status, err := a.xcodeCleaner.Status(ctx)
	if err != nil {
		return err
	}

	if status.EntryCount == 0 {
		fmt.Println("Папок не найдено — нечего чистить.")
		return nil
	}

	fmt.Printf("Найдено папок: %d (%s) в %s\n", status.EntryCount, status.SizeHuman, status.Root)

	confirmed := false
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Удалить все %d папок DerivedData? Действие необратимо.", status.EntryCount),
	}, &confirmed); err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	result := a.xcodeCleaner.Clean(ctx)
	fmt.Printf("Удалено папок: %d\n", result.DeletedCount)
	for _, e := range result.Errs {
		fmt.Println("Ошибка:", e)
	}

	return nil
}

func resolveSelectedTypes(chosen []string, types []enums.CacheType, allOption string) []enums.CacheType {
	for _, c := range chosen {
		if c == allOption {
			return types
		}
	}

	labelToType := make(map[string]enums.CacheType, len(types))
	for _, t := range types {
		labelToType[t.Label()] = t
	}

	selected := make([]enums.CacheType, 0, len(chosen))
	for _, c := range chosen {
		if t, ok := labelToType[c]; ok {
			selected = append(selected, t)
		}
	}

	return selected
}
