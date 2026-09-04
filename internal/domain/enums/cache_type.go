package enums

// CacheType — тип кеша Go, который можно очистить командой `go clean`.
type CacheType int

const (
	CacheTypeBuild CacheType = iota
	CacheTypeMod
	CacheTypeTest
	CacheTypeFuzz
)

// Label — человекочитаемое название для меню и вывода.
func (c CacheType) Label() string {
	switch c {
	case CacheTypeBuild:
		return "build cache (кеш сборки)"
	case CacheTypeMod:
		return "module cache (кеш модулей)"
	case CacheTypeTest:
		return "test cache (кеш тестов)"
	case CacheTypeFuzz:
		return "fuzz cache (кеш фаззинга)"
	default:
		return "неизвестный кеш"
	}
}

// GoCleanFlag — флаг команды `go clean`, соответствующий этому типу кеша.
func (c CacheType) GoCleanFlag() string {
	switch c {
	case CacheTypeBuild:
		return "-cache"
	case CacheTypeMod:
		return "-modcache"
	case CacheTypeTest:
		return "-testcache"
	case CacheTypeFuzz:
		return "-fuzzcache"
	default:
		return ""
	}
}

// AllCacheTypes — все поддерживаемые типы кеша, в порядке отображения в меню.
func AllCacheTypes() []CacheType {
	return []CacheType{CacheTypeBuild, CacheTypeMod, CacheTypeTest, CacheTypeFuzz}
}