package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Iuptec/tupa"
)

var (
	initPgStore   sync.Once
	pgStore       *PostgresStore
	pgInitialized bool
)

func main() {
	server1 := tupa.NewAPIServer(":6969")
	go func() {
		server1.New()
	}()

	server2 := tupa.NewAPIServer(":6968")
	go func() {
		server2.New()
	}()

	pgStore, newInstance := getPostgresStore()
	if newInstance {
		fmt.Println("New instance of PostgresStore created")
	} else {
		fmt.Println("Using existing instance of PostgresStore")
	}

	pgStore.Init()

	server1Manager()
	server1.RegisterRoutes(tupa.GetRoutes())
	// mantem a goroutine principal rodando
	select {}
}

func getPostgresStore() (*PostgresStore, bool) {
	initPgStore.Do(func() {
		pgStore, _ = NewPostgresStore()
		pgStore.Init()
		pgInitialized = true
	})

	// Return the instance and a flag indicating whether it's a new instance
	return pgStore, !pgInitialized
}

func server1Manager() {
	tupa.AddRoutes(nil, server1Routes)
}

func server1Routes() []tupa.RouteInfo {
	return []tupa.RouteInfo{
		{
			Path:    "/clientes/{id}/transacoes",
			Method:  "POST",
			Handler: handleTransacoes,
		},
		{
			Path:   "/hello",
			Method: "GET",
			Handler: func(tc *tupa.TupaContext) error {
				return tc.SendString("Hello world!")
			},
		},
	}
}

func handleTransacoes(ctx *tupa.TupaContext) error {

	var trans transacao = transacao{
		Valor:     1000,
		Tipo:      "c",
		Descricao: "teste",
	}

	id := ctx.Param("id")
	fmt.Printf("ID: %s\n", id)

	db, err := getPostgresStore()
	if err == false {
		fmt.Println("Error getting PostgresStore:", err)
	}

	json := json.NewDecoder(ctx.Request().Body).Decode(&trans)
	if json != nil {
		fmt.Println("Erro decoding JSON:", json)
		return json
	}

	if !db.ClientExists(id) {
		return fmt.Errorf("client with ID %s does not exist", id)
	}

	db.PostInCliente(id, trans)

	fmt.Printf("Valor: %d\n", trans.Valor)
	fmt.Printf("Tipo: %s\n", trans.Tipo)
	fmt.Printf("Descricao: %s\n", trans.Descricao)

	resposta := map[string]int{
		"limite": 10000,
		"saldo":  -9098,
	}

	return tupa.WriteJSONHelper(*ctx.Response(), http.StatusOK, resposta)
}
