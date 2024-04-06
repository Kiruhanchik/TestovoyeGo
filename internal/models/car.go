package models

type CarRequest struct {
	RegNums []string `json:"regNums"`
}

type Car struct {
	Id     *int   `json:"id"`
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  People `json:"owner"`
}

type CarFilter struct {
	RegNum *string       `json:"regNum"`
	Mark   *string       `json:"mark"`
	Model  *string       `json:"model"`
	Year   *int          `json:"year"`
	Owner  *PeopleFilter `json:"owner"`
}

// Package models contains struct for working
