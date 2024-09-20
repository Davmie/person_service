package models

type Person struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Age     int    `json:"age" db:"age"`
	Address string `json:"address" db:"address"`
	Work    string `json:"work" db:"work"`
}
