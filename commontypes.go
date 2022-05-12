package main

import "sync"

type DomainCache struct {
	Cache         map[string][]Rjs
	Notcache      map[string][]Rjs
	cx            sync.RWMutex
	Domains       map[string]*Endings `json:"words"` // what should be keep in here?
	GoogleTlds    []string
	GodaddyTlds   []string
	UniqueTlds    []string
	UniqueTldsRjs []Rjs
	OnlyGoogle    []string
	Xrates        map[string]float64
	GoogleQueries map[string][]Rjs
	gx            sync.RWMutex
	x             sync.RWMutex
}

type Endings struct {
	Tlds map[string]interface{} `json:"words`
	x    sync.RWMutex
}

type Rjs struct {
	Domain string `json:"domain"`
	Tld    string `json:"tld"`
	Price  int64  `json:"price"`
}

type Rj struct {
	Domain   string `json:"domain"`
	Tld      string `json:"tld"`
	Price    int64  `json:"price"`
	Currency string `json:"currency"`
}

type DomainsResponse struct {
	Domains []struct {
		Available  bool   `json:"available"`
		Currency   string `json:"currency"`
		Definitive bool   `json:"definitive"`
		Domain     string `json:"domain"`
		Period     int64  `json:"period"`
		Price      int64  `json:"price"`
	} `json:"domains"`
}

type tldjson []struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
