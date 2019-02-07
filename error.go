package main

import (
	"log"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}