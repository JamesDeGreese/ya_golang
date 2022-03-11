package handlers

type PostJSONRequest struct {
	URL string `json:"url"`
}

type ShortenBatchRequest []struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}
