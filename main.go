package main

import (
	"flag"
	"log"

	"github.com/T1-labRedes/http1"
)

func main() {
	port := flag.Int("port", 8080, "")
	root := flag.String("root", "./www", "")
	flag.Parse()

	srv := http1.NewServer(*port, *root)
	log.Fatal(srv.ListenAndServe())
}
