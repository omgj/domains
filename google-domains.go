package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	domains "cloud.google.com/go/domains/apiv1beta1"
	"github.com/gocolly/colly"
	"github.com/k0kubun/pp"
	domainspb "google.golang.org/genproto/googleapis/cloud/domains/v1beta1"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/genproto/googleapis/type/postaladdress"
)

func TransferParams(domain string) {
	c, e := domains.NewClient(ctx)
	if e != nil {
		log.Fatal(e)
	}
	defer c.Close()

	req := &domainspb.RetrieveTransferParametersRequest{
		DomainName: domain,
		Location:   CloudDomainsLocation,
	}
	resp, e := c.RetrieveTransferParameters(ctx, req)
	if e != nil {
		log.Fatal(e)
	}
	_ = resp.GetTransferParameters()
}

func RegisterParams(domain string) {
	c, err := domains.NewClient(ctx)
	bug(err)
	defer c.Close()
	req := &domainspb.RetrieveRegisterParametersRequest{
		DomainName: domain,
		Location:   CloudDomainsLocation,
	}
	rr, err := c.RetrieveRegisterParameters(ctx, req)
	bug(err)
	a := rr.GetRegisterParameters()
	pp.Print(a.String())
}

// RegisterDomain buys a domain for a user
func RegisterDomain(domain Rj, user string) {
	c, err := domains.NewClient(context.Background())
	bug(err)
	defer c.Close()
	dom := fmt.Sprintf("%s.%s", domain.Domain, domain.Tld)
	req := &domainspb.RegisterDomainRequest{
		Parent: CloudDomainsLocation,
		Registration: &domainspb.Registration{
			DomainName: dom,
			DnsSettings: &domainspb.DnsSettings{
				DnsProvider: &domainspb.DnsSettings_CustomDns_{
					CustomDns: &domainspb.DnsSettings_CustomDns{
						NameServers: CloudDNSNameServers,
					},
				},
			},
			ContactSettings: &domainspb.ContactSettings{
				Privacy: domainspb.ContactPrivacy_PRIVATE_CONTACT_DATA,
				RegistrantContact: &domainspb.ContactSettings_Contact{
					PostalAddress: &postaladdress.PostalAddress{
						RegionCode:         RegistrationCountry,
						PostalCode:         RegistrationPostalCode,
						AdministrativeArea: RegistrationState,
						Locality:           RegistrationCity,
						AddressLines:       AddressLines,
						Recipients:         GoogleRecipients,
					},
					Email:       RegistrationEmail,
					PhoneNumber: RegistrationPhone,
					FaxNumber:   RegistrationFax,
				},
				AdminContact: &domainspb.ContactSettings_Contact{
					PostalAddress: &postaladdress.PostalAddress{
						RegionCode:         RegistrationCountry,
						PostalCode:         RegistrationPostalCode,
						AdministrativeArea: RegistrationState,
						Locality:           RegistrationCity,
						AddressLines:       AddressLines,
						Recipients:         GoogleRecipients,
					},
					Email:       RegistrationEmail,
					PhoneNumber: RegistrationPhone,
					FaxNumber:   RegistrationFax,
				},
				TechnicalContact: &domainspb.ContactSettings_Contact{
					PostalAddress: &postaladdress.PostalAddress{
						RegionCode:         RegistrationCountry,
						PostalCode:         RegistrationPostalCode,
						AdministrativeArea: RegistrationState,
						Locality:           RegistrationCity,
						AddressLines:       AddressLines,
						Recipients:         GoogleRecipients,
					},
					Email:       RegistrationEmail,
					PhoneNumber: RegistrationPhone,
					FaxNumber:   RegistrationFax,
				},
			},
		},
		YearlyPrice: &money.Money{
			CurrencyCode: domain.Currency,
			Units:        domain.Price,
		},
	}
	fmt.Println("Registering")
	op, err := c.RegisterDomain(ctx, req)
	if err != nil {
		a := err.Error()
		if strings.Contains(a, "already registered") {
			fmt.Println("Already registered.")
			return
		}
		if strings.Contains(a, "Wrong yearly price") {
			fmt.Println("Wrong yearly price")
			f := strings.Split(a, " ")
			price := f[len(f)-2]
			fmt.Println("Have: ", domain.Price, "     Need: ", price)
			return
		}
		log.Fatal(err)
	}

	fmt.Println("Waiting")
	rr, err := op.Wait(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	pp.Print(rr)

	regdom(user, dom, "", 11)
}

func DomainSearch(domain string) {
	dc, e := domains.NewClient(ctx)
	if e != nil {
		log.Fatal(e)
	}
	defer dc.Close()

	rq := &domainspb.SearchDomainsRequest{
		Query:    domain,
		Location: "projects/domainsd/locations/global",
	}
	dreg, e := dc.SearchDomains(ctx, rq)
	if e != nil {
		log.Fatal(e)
	}
	qw := make(map[string][]Rjs)
	for _, w := range dreg.GetRegisterParameters() {
		u := w.GetDomainName()
		uu := strings.Split(u, ".")
		domain := uu[0]
		tld := uu[1]
		if len(uu) == 3 {
			tld += "." + uu[2]
		}
		p := Price(w.GetYearlyPrice().GetCurrencyCode(), w.GetYearlyPrice().GetUnits())
		in := false
		fmt.Printf("Domain: %s.%s\tPrice: %d", domain, tld, p)
		for _, po := range qw[domain] {
			if po.Domain == domain && po.Tld == tld && po.Price == p {
				in = true
			}
		}
		if !in {
			fmt.Printf(" adding\n")
			qw[domain] = append(qw[domain], Rjs{domain, tld, p})
			continue
		}
		fmt.Printf(" not adding\n")
	}

	log.Print(qw)
}

func GoogleExtensions() []string {
	c := colly.NewCollector()
	var fs []string
	c.OnHTML("table.nice-table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(an int, e *colly.HTMLElement) {
			var name string
			e.ForEach("td", func(an int, e *colly.HTMLElement) {
				a := e.Text
				if strings.Contains(strings.TrimSpace(a), ".") {
					name = a
				}
				if strings.Contains(a, "$") {
					fs = append(fs, name[1:])
				}
			})
		})
	})
	c.Visit(extensionurl)
	b, e := json.MarshalIndent(fs, "", " ")
	bug(e)
	e = ioutil.WriteFile("google-tlds.json", b, 0644)
	bug(e)
	return fs
}

func googletlds() []string {
	b, e := ioutil.ReadFile("google-tlds.json")
	if os.IsNotExist(e) {
		return GoogleExtensions()
	}
	bug(e)
	var fs []string
	e = json.Unmarshal(b, &fs)
	bug(e)
	return fs
}
