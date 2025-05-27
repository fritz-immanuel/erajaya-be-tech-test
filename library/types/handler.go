package types

//FindAllParams is Parameters helpers for FindAllParams
type FindAllParams struct {
	Page       int
	Size       int
	StatusID   string
	SortBy     string
	SortName   string
	BusinessID string
	CompanyID  string
	DataFinder string
	Outlets    []string
}

// result all
type ResultAll struct {
	Status     string
	Message    string
	StatusCode int
	TotalData  int
	Page       string
	Size       string
	Data       interface{}
}

type Result struct {
	Status     string
	Message    string
	StatusCode int
	Data       interface{}
}
