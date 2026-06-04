package models

type UrlShortenerRequest struct {
	LongUrl string `json:"long_url"`
	UserId  string `json:"user_id"`
}

type UrlShortenerResponse struct {
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
}
