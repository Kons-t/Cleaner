package entities

// IdeaHTTPCacheStatus — сводка по .json-ответам HTTP-клиента IDEA/GoLand в .idea/httpRequests.
type IdeaHTTPCacheStatus struct {
	Root      string
	FileCount int
	SizeHuman string
}

// IdeaHTTPCacheCleanResult — результат удаления .json-ответов HTTP-клиента.
type IdeaHTTPCacheCleanResult struct {
	DeletedCount int
	Errs         []error
}