package model

// ArticleSearchQuery represent article search query parameter
type ArticleSearchQuery struct {
	Keyword    string
	Author     string
	Pagination Pagination
}

// ArticleSearchResult represent article search query result
type ArticleSearchResult struct {
	IDs        []int
	Pagination Pagination
	Total      int
}

// Pagination represent query result pagination parameter
type Pagination struct {
	Offset int
	Limit  int
}
