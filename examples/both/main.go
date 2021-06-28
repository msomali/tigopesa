package main

import "github.com/techcraftt/tigosdk/examples"

func main() {
	err := examples.Server().ListenAndServe()
	if err != nil {
		panic(err)
	}
}
