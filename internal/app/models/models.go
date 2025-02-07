package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURL(id int, short_url string, original_url string) *URL {
	return &URL{
		UUID:        id,
		ShortURL:    short_url,
		OriginalURL: original_url,
	}
}
