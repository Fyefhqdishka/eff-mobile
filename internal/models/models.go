package models

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"releasedate"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
