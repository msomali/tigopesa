package main

import (
	"github.com/joho/godotenv"
	"github.com/techcraftlabs/tigopesa/examples"
	"log"
)

func main() {
	err := godotenv.Load("tigo.env")
	if err != nil {
		log.Printf("error %v\n", err)
		log.Fatal("Error loading .env file")
	}
	err = examples.Server().ListenAndServe()
	if err != nil {
		panic(err)
	}
}
