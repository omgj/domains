package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/siteverification/v1"
)

// SiteVerificationTXT. Set the TXT record returned with
func SiteVerificationTXT(domain string) string {
	vs := SiteVerificationClient()
	fmt.Println("Created new Service Client.")
	req := &siteverification.SiteVerificationWebResourceGettokenRequest{
		Site: &siteverification.SiteVerificationWebResourceGettokenRequestSite{
			Type:       SiteVerificationType,
			Identifier: domain,
		},
		VerificationMethod: SiteVerificationMethod,
	}
	fmt.Println("Fetching TXT data...")
	res, e := vs.WebResource.GetToken(req).Do()
	bug(e)
	token := res.Token
	fmt.Printf("Received: %s\n", token)
	return token
}

// DNSset places the confirmation TXT token and GCP nameservers.
func DNSset(domain, txt string) {
	ctx := context.Background()
	s, e := dns.NewService(ctx)
	bug(e)

	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			{
				Kind:    DnsKind,
				Type:    "TXT",
				Ttl:     300,
				Name:    fmt.Sprintf("%s.", domain),
				Rrdatas: []string{txt},
			},
			{
				Kind:    DnsKind,
				Type:    "A",
				Ttl:     300,
				Name:    fmt.Sprintf("%s.", domain),
				Rrdatas: DnsAs,
			},
			{
				Kind:    DnsKind,
				Type:    "AAAA",
				Ttl:     300,
				Name:    fmt.Sprintf("%s.", domain),
				Rrdatas: DnsAAAAs,
			},
		},
	}

	cc, e := s.Changes.Create(GCPProject, GCPManagedZone, change).Do()
	bug(e)

	fmt.Println(cc.HTTPStatusCode)
	fmt.Println(cc.Status)

}

// TryVerify attempts to verify the domain. Set TXT record first
func TryVerify(domain string) {
	vs := SiteVerificationClient()
	sr := &siteverification.SiteVerificationWebResourceResource{
		Site: &siteverification.SiteVerificationWebResourceResourceSite{
			Type:       SiteVerificationType,
			Identifier: domain,
		},
	}
	fmt.Printf("Attempting Verification of %s\n", domain)
	insertcall, e := vs.WebResource.Insert(SiteVerificationType, sr).Do()
	bug(e)
	if insertcall.HTTPStatusCode != 200 {
		fmt.Println("Verification Failed.")
		log.Print(insertcall)
	}
}

func SiteVerificationClient() *siteverification.Service {
	ctx := context.Background()
	// Set up credentials with user impersonation.
	jsonKey, _ := ioutil.ReadFile("/Users/oliver/master/gkey.json")
	jwt, e := google.JWTConfigFromJSON(jsonKey, SiteVerificationJWT)
	bug(e)
	jwt.Subject = SuperEmail
	tokenSource := jwt.TokenSource(ctx)
	vs, e := siteverification.NewService(ctx, option.WithTokenSource(tokenSource))
	bug(e)
	return vs
}

// ListVerified lists verified domains of SuperEmail & siliconuxx
func ListVerified() {
	jsonKey, _ := ioutil.ReadFile("/Users/oliver/master/gkey.json")
	jwt, e := google.JWTConfigFromJSON(jsonKey, "https://www.googleapis.com/auth/siteverification")
	bug(e)
	jwt.Subject = SuperEmail
	tokenSource := jwt.TokenSource(ctx)
	verificationService, e := siteverification.NewService(ctx, option.WithTokenSource(tokenSource))
	bug(e)

	ds, e := verificationService.WebResource.List().Do()
	bug(e)
	for _, a := range ds.Items {
		fmt.Println(a.Site.Identifier)
	}
}
