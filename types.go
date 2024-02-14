package main

type transacao struct {
	Valor     int64  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}
