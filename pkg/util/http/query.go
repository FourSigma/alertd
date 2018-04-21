package http

type Filter struct {
	Offset int64  `schema:"offset"`
	Limit  int64  `schema:"limit"`
	SortBy string `schema:"sortby"`
	Asc    bool   `schema:"asc"`
}
