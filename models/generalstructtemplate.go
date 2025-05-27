package models

type IDNameTemplate struct {
	ID   int    `json:"ID" db:"id"`
	Name string `json:"Name" db:"name"`
}

type StringIDNameTemplate struct {
	ID   string `json:"ID" db:"id"`
	Name string `json:"Name" db:"name"`
}

type StringIDNameCodeTemplate struct {
	ID   string `json:"ID" db:"id"`
	Name string `json:"Name" db:"name"`
	Code string `json:"Code" db:"code"`
}
