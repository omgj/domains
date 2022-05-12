package main

import "log"

func bug(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
