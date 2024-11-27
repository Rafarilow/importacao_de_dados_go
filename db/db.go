package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Necessário para usar o driver MySQL
)

var db *sql.DB

// Conecta ao banco de dados MySQL
func ConnectDB(user, password, host, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Testa a conexão
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Fecha a conexão com o banco de dados
func Close() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

// Função para inserir contato no banco de dados
func InsertContact(contact map[string]string) error {
	_, err := db.Exec("INSERT INTO contatos (nome, email, telefone) VALUES (?, ?, ?)",
		contact["Nome"], contact["Email"], contact["Telefone"])
	return err
}
