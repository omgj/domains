package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/pp"
)

type GdAgreement []struct {
	AgreementKey string `json:"agreementKey"`
	Content      string `json:"content"`
	Title        string `json:"title"`
	URL          string `json:"url"`
}

func GdAgreements(domain Rj) GdAgreement {
	c, e := http.NewRequest("GET", fmt.Sprintf("https://api.godaddy.com/v1/domains/agreements?tlds=%s&privacy=true", domain.Tld), nil)
	bug(e)
	c.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	c.Header.Set("X-Market-Id", "en-US")
	c.Header.Set("Accept", "application/json")
	a, e := http.DefaultClient.Do(c)
	bug(e)
	defer a.Body.Close()
	fmt.Println(a.StatusCode)
	if a.StatusCode == 200 {
		var gd GdAgreement
		json.NewDecoder(a.Body).Decode(&gd)
		return gd
	}
	return GdAgreement{}
}

type GodaddySchema struct {
	Consent struct {
		AgreedAt      string   `json:"agreedAt"`
		AgreedBy      string   `json:"agreedBy"`
		AgreementKeys []string `json:"agreementKeys"`
	} `json:"consent"`
	ContactAdmin struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactAdmin"`
	ContactBilling struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactBilling"`
	ContactRegistrant struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactRegistrant"`
	ContactTech struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactTech"`
	Domain      string   `json:"domain"`
	NameServers []string `json:"nameServers"`
	Period      int      `json:"period"`
	Privacy     bool     `json:"privacy"`
	RenewAuto   bool     `json:"renewAuto"`
}

func Gdvalidate(domain Rj, gd GdAgreement) {
	var a GodaddySchema
	var keys []string
	for _, b := range gd {
		keys = append(keys, b.AgreementKey)
	}
	a.Consent.AgreementKeys = keys
	a.Consent.AgreedAt = time.Now().UTC().Format(time.RFC3339)
	a.Consent.AgreedBy = RegistrationIP
	a.ContactAdmin.AddressMailing.Address1 = RegistrationAddress1
	a.ContactAdmin.AddressMailing.Address2 = RegistrationAddress2
	a.ContactAdmin.AddressMailing.City = RegistrationCity
	a.ContactAdmin.AddressMailing.Country = RegistrationCountry
	a.ContactAdmin.AddressMailing.PostalCode = RegistrationPostalCode
	a.ContactAdmin.AddressMailing.State = RegistrationState
	a.ContactAdmin.Email = RegistrationEmail
	a.ContactAdmin.Fax = RegistrationFax
	a.ContactAdmin.JobTitle = RegistrationJobTitle
	a.ContactAdmin.NameFirst = RegistrationNameFirst
	a.ContactAdmin.NameLast = RegistrationNameLast
	a.ContactAdmin.NameMiddle = RegistrationNameMiddle
	a.ContactAdmin.Organization = RegistrationOrganization
	a.ContactAdmin.Phone = RegistrationPhone
	a.ContactBilling.AddressMailing.Address1 = RegistrationAddress1
	a.ContactBilling.AddressMailing.Address2 = RegistrationAddress2
	a.ContactBilling.AddressMailing.City = RegistrationCity
	a.ContactBilling.AddressMailing.Country = RegistrationCountry
	a.ContactBilling.AddressMailing.PostalCode = RegistrationPostalCode
	a.ContactBilling.AddressMailing.State = RegistrationState
	a.ContactBilling.Email = RegistrationEmail
	a.ContactBilling.Fax = RegistrationFax
	a.ContactBilling.JobTitle = RegistrationJobTitle
	a.ContactBilling.NameFirst = RegistrationNameFirst
	a.ContactBilling.NameLast = RegistrationNameLast
	a.ContactBilling.NameMiddle = RegistrationNameMiddle
	a.ContactBilling.Organization = RegistrationOrganization
	a.ContactBilling.Phone = RegistrationPhone
	a.ContactRegistrant.AddressMailing.Address1 = RegistrationAddress1
	a.ContactRegistrant.AddressMailing.Address2 = RegistrationAddress2
	a.ContactRegistrant.AddressMailing.City = RegistrationCity
	a.ContactRegistrant.AddressMailing.Country = RegistrationCountry
	a.ContactRegistrant.AddressMailing.PostalCode = RegistrationPostalCode
	a.ContactRegistrant.AddressMailing.State = RegistrationState
	a.ContactRegistrant.Email = RegistrationEmail
	a.ContactRegistrant.Fax = RegistrationFax
	a.ContactRegistrant.JobTitle = RegistrationJobTitle
	a.ContactRegistrant.NameFirst = RegistrationNameFirst
	a.ContactRegistrant.NameLast = RegistrationNameLast
	a.ContactRegistrant.NameMiddle = RegistrationNameMiddle
	a.ContactRegistrant.Organization = RegistrationOrganization
	a.ContactRegistrant.Phone = RegistrationPhone
	a.ContactTech.AddressMailing.Address1 = RegistrationAddress1
	a.ContactTech.AddressMailing.Address2 = RegistrationAddress2
	a.ContactTech.AddressMailing.City = RegistrationCity
	a.ContactTech.AddressMailing.Country = RegistrationCountry
	a.ContactTech.AddressMailing.PostalCode = RegistrationPostalCode
	a.ContactTech.AddressMailing.State = RegistrationState
	a.ContactTech.Email = RegistrationEmail
	a.ContactTech.Fax = RegistrationFax
	a.ContactTech.JobTitle = RegistrationJobTitle
	a.ContactTech.NameFirst = RegistrationNameFirst
	a.ContactTech.NameLast = RegistrationNameLast
	a.ContactTech.NameMiddle = RegistrationNameMiddle
	a.ContactTech.Organization = RegistrationOrganization
	a.ContactTech.Phone = RegistrationPhone
	a.Domain = fmt.Sprintf("%s.%s", domain.Domain, domain.Tld)
	a.NameServers = CloudDNSNameServers
	a.Period = 1
	a.Privacy = true
	a.RenewAuto = true

	b, e := json.Marshal(a)
	bug(e)
	c, e := http.NewRequest("POST", "https://api.godaddy.com/v1/domains/purchase/validate", bytes.NewBuffer(b))
	bug(e)
	c.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	c.Header.Set("Content-Type", "application/json")
	c.Header.Set("Accept", "application/json")
	d, e := http.DefaultClient.Do(c)
	bug(e)
	defer d.Body.Close()
	fmt.Println("Status: ", d.StatusCode)
	if d.StatusCode == 200 {
		fmt.Println("Yes.")
		return
	}
	if d.StatusCode == 400 {
		fmt.Println("malformed request")
		var dd GodaddySchemaErr
		json.NewDecoder(d.Body).Decode(&dd)
		pp.Print(dd)
	}
	if d.StatusCode == 401 {
		fmt.Println("authentication info not sent or invalid")
		var dd GodaddySchemaErr
		json.NewDecoder(d.Body).Decode(&dd)
		pp.Print(dd)
	}
	if d.StatusCode == 429 {
		var w GodaddyErr
		e = json.NewDecoder(d.Body).Decode(&w)
		bug(e)
		pp.Print(w)
		return
	}
	var w GodaddySchemaErr
	e = json.NewDecoder(d.Body).Decode(&w)
	bug(e)
	pp.Print(w)
}

func GodaddyExtensions() tldjson {
	r, e := http.NewRequest("GET", "https://api.godaddy.com/v1/domains/tlds", nil)
	bug(e)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	res, e := http.DefaultClient.Do(r)
	bug(e)
	defer res.Body.Close()
	var tlds tldjson
	json.NewDecoder(res.Body).Decode(&tlds)
	return tlds
}

func godaddytlds() []string {
	var tlds []string
	b, e := ioutil.ReadFile("godaddy-tlds.json")
	if os.IsNotExist(e) {
		for _, a := range GodaddyExtensions() {
			tlds = append(tlds, a.Name)
		}
		w, e := json.MarshalIndent(tlds, "", " ")
		bug(e)
		e = ioutil.WriteFile("godaddy-tlds.json", w, 0644)
		bug(e)
		return tlds
	}
	bug(e)
	bug(json.Unmarshal(b, &tlds))
	return tlds
}

func gdget(domain string, tld string) Rj {
	r, e := http.NewRequest("GET", fmt.Sprintf("https://api.godaddy.com/v1/domains/available?domain=%s.%s&checkType=FULL", domain, tld), nil)
	if e != nil {
		log.Fatal(e)
	}
	r.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	q, e := http.DefaultClient.Do(r)
	if e != nil {
		log.Fatal(e)
	}
	defer q.Body.Close()

	if q.StatusCode == 429 {
		var er GodaddyErr
		e := json.NewDecoder(q.Body).Decode(&er)
		if e != nil {
			log.Fatal(e)
		}
		return Rj{Price: int64(er.RetryAfterSec)}
	}

	if q.StatusCode == 200 {
		var er GodaddyGet
		e := json.NewDecoder(q.Body).Decode(&er)
		if e != nil {
			log.Fatal(e)
		}
		if er.Price > 0 {
			er.Price /= 1000000
		}
		fmt.Printf("Godaddy GET: %s\t%d\t%s %d\n", fmt.Sprintf("%s.%s", domain, tld), q.StatusCode, er.Currency, er.Price)
		return Rj{domain, tld, er.Price, er.Currency}
	}
	return Rj{}
}

// DOMAIN INFORMATION

type GodaddyDomainInfo struct {
	AuthCode     string `json:"authCode"`
	ContactAdmin struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactAdmin"`
	ContactBilling struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactBilling"`
	ContactRegistrant struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactRegistrant"`
	ContactTech struct {
		AddressMailing struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			State      string `json:"state"`
		} `json:"addressMailing"`
		Email        string `json:"email"`
		Fax          string `json:"fax"`
		JobTitle     string `json:"jobTitle"`
		NameFirst    string `json:"nameFirst"`
		NameLast     string `json:"nameLast"`
		NameMiddle   string `json:"nameMiddle"`
		Organization string `json:"organization"`
		Phone        string `json:"phone"`
	} `json:"contactTech"`
	CreatedAt              time.Time `json:"createdAt"`
	DeletedAt              time.Time `json:"deletedAt"`
	TransferAwayEligibleAt time.Time `json:"transferAwayEligibleAt"`
	Domain                 string    `json:"domain"`
	DomainID               int       `json:"domainId"`
	ExpirationProtected    bool      `json:"expirationProtected"`
	Expires                time.Time `json:"expires"`
	ExposeWhois            bool      `json:"exposeWhois"`
	HoldRegistrar          bool      `json:"holdRegistrar"`
	Locked                 bool      `json:"locked"`
	NameServers            []string  `json:"nameServers"`
	Privacy                bool      `json:"privacy"`
	RenewAuto              bool      `json:"renewAuto"`
	RenewDeadline          time.Time `json:"renewDeadline"`
	Status                 string    `json:"status"`
	SubaccountID           string    `json:"subaccountId"`
	TransferProtected      bool      `json:"transferProtected"`
	Verifications          struct {
		DomainName struct {
			Status string `json:"status"`
		} `json:"domainName"`
		RealName struct {
			Status string `json:"status"`
		} `json:"realName"`
	} `json:"verifications"`
}

func GodaddyDomainInfos(domain string) {
	r, e := http.NewRequest("GET", fmt.Sprintf("https://api.godaddy.com/v1/domains/%s", domain), nil)
	bug(e)
	r.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	q, e := http.DefaultClient.Do(r)
	if e != nil {
		log.Fatal(e)
	}
	defer q.Body.Close()
	fmt.Printf("Godaddy GET: %s\t%d\n", domain, q.StatusCode)

	if q.StatusCode == 200 {
		var er GodaddyDomainInfo
		e := json.NewDecoder(q.Body).Decode(&er)
		bug(e)
		pp.Print(er)
		return
	}

	var gd GodaddySchemaErr
	e = json.NewDecoder(q.Body).Decode(&gd)
	bug(e)
	pp.Print(gd)
}

// Crawl will back off when asked too.
func Crawl() {
	tlds := Dcache.GodaddyTlds
	for _, v := range alphabet {
		for _, a := range alphabet {
			for _, z := range alphabet {
				var res []Rj
				dom := fmt.Sprintf("%s%s%s", string(v), string(a), string(z))
				for _, b := range tlds {
					c := gdget(dom, b)
					if c.Domain == "" {
						fmt.Println("Max out. Sleeping ", c.Price)
						time.Sleep(time.Second * time.Duration(c.Price))
						res = append(res, gdget(dom, b))
						continue
					}
					res = append(res, c)
				}
				fmt.Println("writing file")
				t, e := json.MarshalIndent(res, "", " ")
				if e != nil {
					log.Fatal(e)
				}
				e = ioutil.WriteFile(fmt.Sprintf("%s.json", dom), t, 0644)
				if e != nil {
					log.Fatal(e)
				}
			}
		}
		return
	}
}

type ValInfo struct {
	ID         string   `json:"id"`
	Models     struct{} `json:"models"`
	Properties struct{} `json:"properties"`
	Required   []string `json:"required"`
}

func Gdschema(tlds string) {
	c, e := http.NewRequest("GET", fmt.Sprintf("https://api.godaddy.com/v1/domains/purchase/schema/%s", tlds), nil)
	bug(e)
	c.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GodaddyKey, GodaddySecret))
	c.Header.Set("Content-Type", "application/json")
	c.Header.Set("Accept", "application/json")
	a, e := http.DefaultClient.Do(c)
	bug(e)
	defer a.Body.Close()
	if a.StatusCode == 200 {
		var v ValInfo
		e := json.NewDecoder(a.Body).Decode(&v)
		bug(e)
		pp.Print(v)
		return
	}
	if a.StatusCode == 429 {
		var w GodaddyErr
		e = json.NewDecoder(a.Body).Decode(&w)
		bug(e)
		pp.Print(w)
		return
	}
	var w GodaddySchemaErr
	e = json.NewDecoder(a.Body).Decode(&w)
	bug(e)
	pp.Print(w)
}

type GodaddySchemaErr struct {
	Code   string `json:"code"`
	Fields []struct {
		Code        string `json:"code"`
		Message     string `json:"message"`
		Path        string `json:"path"`
		PathRelated string `json:"pathRelated"`
	} `json:"fields"`
	Message string `json:"message"`
}

type GodaddyErr struct {
	Code   string `json:"code"`
	Fields []struct {
		Code        string `json:"code"`
		Message     string `json:"message"`
		Path        string `json:"path"`
		PathRelated string `json:"pathRelated"`
	} `json:"fields"`
	Message       string `json:"message"`
	RetryAfterSec int    `json:"retryAfterSec"`
}

type GodaddyGet struct {
	Available  bool   `json:"available"`
	Currency   string `json:"currency"`
	Definitive bool   `json:"definitive"`
	Domain     string `json:"domain"`
	Period     int64  `json:"period"`
	Price      int64  `json:"price"`
}
