package tmdb

type Movie struct {
	Id int `json:"id"`
}

type LatestResponse struct {
	Results []Movie `json:"results"`
}
