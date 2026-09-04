// Package paths собирает файловые пути, специфичные для конкретной машины/пользователя.
// Каждый путь можно переопределить через переменную окружения (см. .env.example в корне
// репозитория) — без этого используется дефолт, вычисленный от домашнего каталога пользователя.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var homeDir = resolveHomeDir()

// WorkDir — корень поиска файлов для очистки (например, .idea/httpRequests с ответами
// HTTP-клиента). "Go_projects/Work" — не стандартный путь ни на одной системе, поэтому
// по умолчанию берётся весь домашний каталог; сузить область поиска до своих реальных
// рабочих проектов можно через CLEANER_WORK_DIR.
var WorkDir = envOr("CLEANER_WORK_DIR", homeDir)

// LogsDir — каталог логов приложений текущего пользователя. Переопределяется CLEANER_LOGS_DIR.
var LogsDir = envOr("CLEANER_LOGS_DIR", filepath.Join(homeDir, "Library", "Logs"))

// DSStoreRoot — корень поиска мусорных .DS_Store файлов Finder. Переопределяется CLEANER_DS_STORE_ROOT.
var DSStoreRoot = envOr("CLEANER_DS_STORE_ROOT", homeDir)

// DavinciCacheDir — стандартный каталог render cache DaVinci Resolve на macOS
// (используется, если в настройках Resolve не задан кастомный путь кеша).
// Переопределяется CLEANER_DAVINCI_CACHE_DIR.
var DavinciCacheDir = envOr("CLEANER_DAVINCI_CACHE_DIR", filepath.Join(homeDir, "Movies", "CacheClip"))

// DavinciAppPath — путь установки DaVinci Resolve, используется для проверки,
// установлено ли приложение на этом устройстве. Переопределяется CLEANER_DAVINCI_APP_PATH.
var DavinciAppPath = envOr("CLEANER_DAVINCI_APP_PATH", "/Applications/DaVinci Resolve")

// XcodeDerivedDataDir — каталог DerivedData Xcode: build-артефакты и индексы,
// сгруппированные по подпапкам на проект. Переопределяется CLEANER_XCODE_DERIVED_DATA_DIR.
var XcodeDerivedDataDir = envOr(
	"CLEANER_XCODE_DERIVED_DATA_DIR",
	filepath.Join(homeDir, "Library", "Developer", "Xcode", "DerivedData"),
)

// resolveHomeDir определяет домашний каталог пользователя. Если это не удаётся —
// программа завершается с понятной ошибкой вместо использования вшитого чужого пути.
func resolveHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil || dir == "" {
		fmt.Fprintln(os.Stderr, "cleaner: не удалось определить домашний каталог пользователя (os.UserHomeDir).")
		fmt.Fprintln(os.Stderr, "Задайте нужные пути вручную через переменные окружения — см. .env.example.")
		os.Exit(1)
	}

	return dir
}

var loadEnvFileOnce sync.Once

// envOr возвращает значение переменной окружения key, если оно задано и непусто, иначе fallback.
// Перед первым чтением подгружает .env из текущей рабочей директории, если он есть
// (это не влияет на уже экспортированные переменные окружения — они имеют приоритет).
func envOr(key, fallback string) string {
	loadEnvFileOnce.Do(func() {
		_ = godotenv.Load()
	})

	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}