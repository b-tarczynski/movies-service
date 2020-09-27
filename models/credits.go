package models

type Credit struct {
	_       struct{} `pg:",discard_unknown_columns"`
	Id      int      `json:"id"`
	MovieId int      `json:"-"`
	Cast    []*Cast  `json:"cast" pg:"fk:credit_id"`
	Crew    []*Crew  `json:"crew" pg:"fk:credit_id"`
}

type Cast struct {
	CastId      int    `json:"cast_id"`
	Id          int    `json:"id"`
	CreditId    int    `json:"-"`
	Gender      int    `json:"gender"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	Order       int    `json:"order"`
	ProfilePath string `json:"profile_path"`
}

type Crew struct {
	Id          int    `json:"id"`
	CreditId    int    `json:"-"`
	Department  string `json:"department"`
	Gender      int    `json:"gender"`
	Job         string `json:"job"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}
