package models

type Unit struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	UnitID   int    `json:"unit_id"`
	Unit     Unit   `json:"unit,omitempty"`
}
