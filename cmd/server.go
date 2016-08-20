package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yhanada/apimock"
)

var (
	port  uint
	root  string
	check bool
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.UintVar(&port, "port", 8080, "server port")
	flag.StringVar(&root, "root", wd, "Document Root")
	flag.BoolVar(&check, "check", false, "Check actions")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  %s [OPTIONS]\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	ctx, err := apimock.NewContext(root)
	if err != nil {
		panic(err)
	}

	if check {
		ctx.CheckActions()
		return
	}

	http.HandleFunc("/", ctx.GetHandlerFunc())

	log.Println("Start API Mock Server. Port:", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("Failed to start Server:", err)
	}
}
