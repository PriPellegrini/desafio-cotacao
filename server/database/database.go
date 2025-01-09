package database

import (
	"context"
	"database/sql"

	"github.com/PriPellegrini/desafio-cotacao/cotacao"
)

type CotacaoDB struct {
	DB *sql.DB
}

func NewCotacaoDB(db *sql.DB) *CotacaoDB {
	return &CotacaoDB{DB: db}
}

func (d CotacaoDB) Insert(ctx context.Context, cotacao cotacao.Cotacao) error {

	insert := "INSERT INTO cotacoes (code, codein, name, high, low, varbid, pctchange, bid, ask, timestamp, createDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := d.DB.ExecContext(ctx,insert, cotacao.Usdbrl.Code, cotacao.Usdbrl.Codein, cotacao.Usdbrl.Name, cotacao.Usdbrl.High, cotacao.Usdbrl.Low, cotacao.Usdbrl.VarBid, cotacao.Usdbrl.PctChange, cotacao.Usdbrl.Bid, cotacao.Usdbrl.Ask, cotacao.Usdbrl.Timestamp, cotacao.Usdbrl.CreateDate)
	if err != nil {
		panic(err)
	}
	return nil
}
