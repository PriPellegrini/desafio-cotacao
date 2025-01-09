package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro de timeout")
	}
	defer res.Body.Close()

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	resCotacao, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write([]byte("Dólar: " + string(resCotacao)))
	if err != nil {
		panic(err)
	}
	fmt.Println("arquivo criado com sucesso")

	file.Close()

}
