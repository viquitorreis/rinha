package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", connStr)
	// se o db não existe fica eternamente carregando a api
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (p *PostgresStore) Init() error {
	if err := p.CreateSaldoTable(); err != nil {
		return fmt.Errorf("error creating saldo table: %v", err)
	}
	if err := p.CreateClienteTable(); err != nil {
		return fmt.Errorf("error creating cliente table: %v", err)
	}
	if err := p.CreateTransacoesTable(); err != nil {
		return fmt.Errorf("error creating transacoes table: %v", err)
	}
	if err := p.startWith5Clientes(); err != nil {
		return fmt.Errorf("error starting with 5 clientes: %v", err)
	}
	return nil
}

func (p *PostgresStore) CreateClienteTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS clientes (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			saldo INT,
            FOREIGN KEY (saldo) REFERENCES saldo(id)
		)
	`

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) CreateSaldoTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS saldo (
			id SERIAL PRIMARY KEY,
			total INT,
			data_extrato DATE,
			limite INT
		)
	`

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) CreateTransacoesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS transacoes (
            id SERIAL PRIMARY KEY,
            cliente_id INT,
            valor INT,
            tipo VARCHAR(1),
            descricao VARCHAR(100),
            realizada_em DATE,
            FOREIGN KEY (cliente_id) REFERENCES clientes(id)
        )
    `

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresStore) PostInCliente(clienteId string, transac transacao) {
	valor := transac.Valor
	tipo := transac.Tipo
	descricao := transac.Descricao
	clientId := clienteId
	query := fmt.Sprintf("INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES ('%s', %d, '%s', '%s')", clientId, valor, tipo, descricao)
	_, err := p.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PostgresStore) startWith5Clientes() error {
	_, err := p.db.Exec("INSERT INTO clientes (nome, saldo) VALUES ('cliente1', 1000)")
	if err != nil {
		return fmt.Errorf("error inserting cliente1: %v", err)
	}

	_, err = p.db.Exec("INSERT INTO clientes (nome, saldo) VALUES ('cliente2', 1000)")
	if err != nil {
		return fmt.Errorf("error inserting cliente2: %v", err)
	}

	_, err = p.db.Exec("INSERT INTO clientes (nome, saldo) VALUES ('cliente3', 1000)")
	if err != nil {
		return fmt.Errorf("error inserting cliente3: %v", err)
	}

	_, err = p.db.Exec("INSERT INTO clientes (nome, saldo) VALUES ('cliente4', 1000)")
	if err != nil {
		return fmt.Errorf("error inserting cliente4: %v", err)
	}

	_, err = p.db.Exec("INSERT INTO clientes (nome, saldo) VALUES ('cliente5', 1000)")
	if err != nil {
		return fmt.Errorf("error inserting cliente5: %v", err)
	}

	return nil
}

func (p *PostgresStore) ClientExists(clientID string) bool {
	query := "SELECT COUNT(*) FROM clientes WHERE id = $1"
	var count int
	err := p.db.QueryRow(query, clientID).Scan(&count)
	if err != nil {
		log.Println("Error checking if client exists:", err)
		return false
	}
	return count > 0
}
