package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Web(w http.ResponseWriter, r *http.Request) {
	b, e := ioutil.ReadFile("index.html")
	if e != nil {
		log.Fatal(e)
	}
	w.Write(b)
}

func GetTlds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rjs []Rjs
	for _, a := range Dcache.GodaddyTlds {
		rjs = append(rjs, Rjs{"", a, 0})
	}
	json.NewEncoder(w).Encode(rjs)
}

func GetCache(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query()
	domain := a["q"][0]

	if domain == "" {
		fmt.Println("Bad cache get")
		var r []Rjs
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(r)
		return
	}
	fmt.Println("Cache return for: ", domain)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Dcache.ReadCache(domain))
}

func TryRegister(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&m)
	log.Println(m)
}
func MakeCache() {

}
