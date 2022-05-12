package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

func NewAccount(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query()
	uid := a["uid"][0]
	email := a["email"][0]
	name := a["name"][0]
	ip := r.RemoteAddr
	// us := r.UserAgent()
	cid, aid := StripeCustomerAccount(name, email)

	ctx := context.Background()
	c, err := firestore.NewClient(ctx, "domainsd")
	if err != nil {
		log.Fatal(err)
	}
	c.Collection("users").Doc(uid).Set(context.Background(), map[string]interface{}{
		"cid": cid,
		"aid": aid,
		"ips": ip,
	})

}

func regdom(user, domain, pi string, amount int64) {
	ctx := context.Background()
	c, err := firestore.NewClient(ctx, "domainsd")
	if err != nil {
		log.Fatal(err)
	}
	gg, e := c.Collection("registrations").Doc(user).Set(ctx, map[string]interface{}{
		"domain": domain,
		"user":   user,
		"amount": amount,
	})
	if e != nil {
		log.Println(e)
		log.Println(gg)
	}
}
