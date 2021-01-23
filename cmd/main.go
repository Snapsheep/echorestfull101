package main

import (
	"nevergo"
	_ "nevergo/docs"

	_ "github.com/lib/pq"
)

func main() {
	nevergo.StartServices()
}
