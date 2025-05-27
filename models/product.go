package models

import (
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
)

type ProductBulk struct {
	ID          string    `json:"ID" db:"id"`
	Name        string    `json:"Name" db:"name"`
	Price       float64   `json:"Price" db:"price"`
	Description string    `json:"Description" db:"description"`
	Quantity    int       `json:"Quantity" db:"quantity"`
	StatusID    string    `json:"StatusID" db:"status_id"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`

	StatusName string `json:"StatusName" db:"status_name"`
}

type Product struct {
	ID          string    `json:"ID" db:"id"`
	Name        string    `json:"Name" db:"name"`
	Price       float64   `json:"Price" db:"price"`
	Description string    `json:"Description" db:"description"`
	Quantity    int       `json:"Quantity" db:"quantity"`
	StatusID    string    `json:"StatusID" db:"status_id"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`

	Status Status `json:"Status"`
}

type FindAllProductParams struct {
	FindAllParams types.FindAllParams
}
