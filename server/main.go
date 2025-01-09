package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/PriPellegrini/desafio-cotacao/cotacao"
	"github.com/PriPellegrini/desafio-cotacao/database"
	_ "modernc.org/sqlite"
)

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro de timeout")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var c cotacao.Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs("./cotacoes.db")
	if err != nil {
		log.Fatal("Erro ao obter o caminho absoluto do banco de dados:", err)
	}

	fmt.Println("Caminho absoluto do banco de dados:", absPath)

	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		log.Fatal("Erro ao abrir o banco de dados: ", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varbid TEXT,
		pctchange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		createDate TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Erro ao criar a tabela:", err)
	}

	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='cotacoes';")
	var tableName string
	err = row.Scan(&tableName)
	if err != nil {
		log.Fatal("Erro ao verificar a tabela:", err)
	}

	cotacaoDB := database.NewCotacaoDB(db)

	ctxInsert, cancelInsert := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancelInsert()

	err = cotacaoDB.Insert(ctxInsert, c)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(c.Usdbrl.Bid))

	fmt.Println(string(body))

}
