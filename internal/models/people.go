package models

type People struct {
	Id         int    `json:"Id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

// Package models contains struct for working
