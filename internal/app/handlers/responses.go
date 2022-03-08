package handlers

type PostJSONResponse struct {
	Result string `json:"result"`
}

type UserLinkItem struct {
	SortURL     string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
