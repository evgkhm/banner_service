package entity

type Product struct {
	ID         int               `json:"id" example:"1"`
	Properties map[string]string `json:"properties" example:"{\"штрихкод\":\"12345678\"}"`
}
