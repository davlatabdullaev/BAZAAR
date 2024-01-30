package models

type Response struct {
	StatusCode  int
	Description string
	Data        interface{}
}
