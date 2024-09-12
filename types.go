package main

type ShortenURLRequestBody struct {
	Url string `json:"url" form:"url" binding:"required"`
}

type ShortenURLResponse struct {
	MiniURL string `json:"miniurl"`
	Code    int    `json:"code"`
}

type Error struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}
