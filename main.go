package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("bad port")
	}
	a := App{}
	a.Initialize(os.Getenv("DATABASE_URL"))
	a.Run(port)
}