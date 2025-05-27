package models

import (
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
)

type UserActionBulk struct {
	ID          string    `json:"ID" db:"id"`
	UserID      string    `json:"UserID" db:"user_id"`
	TableName   string    `json:"TableName" db:"table_name"`
	Action      string    `json:"Action" db:"action"`
	ActionValue int       `json:"ActionValue" db:"action_value"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`

	StatusName string `json:"StatusName" db:"status_name"`
	RefID      string `json:"RefID" db:"ref_id"`
}

type UserAction struct {
	ID          string    `json:"ID" db:"id"`
	UserID      string    `json:"UserID" db:"user_id"`
	TableName   string    `json:"TableName" db:"table_name"`
	Action      string    `json:"Action" db:"action"`
	ActionValue int       `json:"ActionValue" db:"action_value"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`

	StatusName string `json:"StatusName"`
	RefID      string `json:"RefID" db:"ref_id"`
}

type FindAllActionHistory struct {
	FindAllParams    types.FindAllParams
	UsingStatusTable int
	UserID           string
	RefID            string
	ModuleName       string
	TableName        string
	PackageName      string
	GroupBy          string
}
