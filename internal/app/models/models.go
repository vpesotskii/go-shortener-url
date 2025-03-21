package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url"`
}

type URL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURL(id int, shortURL string, originalURL string) *URL {
	return &URL{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
