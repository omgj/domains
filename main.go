package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
)

var (
	fire      *firebase.App
	Dcache    *DomainCache
	ctx       = context.Background()
	backoff   = 0
	popgodad  = []string{}
	chinese   = "xn--6frz82g"
	chnesealt = "移动"
	alphabet  = "abcdefghijklmnopqrstuvwxyz"

	// Currency Rates - Domains are priced mostly in USD, will need to convert.
	XratesKEY = ""
	XratesAPI = "http://data.fixer.io/api/latest?access_key=%s"

	// Stripe
	StripeKEY  = ""
	RETURN_URL = ""

	// Godaddy - removed secrets
	GodaddySecret            = ""
	GodaddyKey               = ""
	RegistrationJobTitle     = "" // Boss
	RegistrationNameFirst    = ""
	RegistrationNameLast     = ""
	RegistrationNameMiddle   = ""
	RegistrationOrganization = "" // Org
	GoogleRegistrationPhone  = ""
	RegistrationAddress1     = "" // Unit Number
	RegistrationAddress2     = "" // Street Address
	RegistrationIP           = ""

	// Google - removed secrets
	extensionurl             = `https://support.google.com/domains/answer/6010092?hl=en#zippy=%2Cprice-by-domain-ending`
	GoogleDomainsNameServers = []string{`ns-cloud-b1.googledomains.com.`, `ns-cloud-b2.googledomains.com.`, `ns-cloud-b3.googledomains.com.`, `ns-cloud-b4.googledomains.com`}
	CloudDNSNameServers      = []string{`ns-cloud-c1.googledomains.com.`, `ns-cloud-c2.googledomains.com.`, `ns-cloud-c3.googledomains.com.`, `ns-cloud-c4.googledomains.com`}
	CloudDomainsLocation     = `` // projects/PROJECT_NAME/locations/global
	IAMServiceAccount        = `` // SERVICE_ACCOUNT_NAME@PROJECT_NAME.iam.gserviceaccount.com
	SuperEmail               = `` // Google Super user Email for Sire Verification i.e someone@customdomain.tld
	SiteVerificationJWT      = `https://www.googleapis.com/auth/siteverification`
	SiteVerificationType     = `INET_DOMAIN`
	SiteVerificationMethod   = `DNS_TXT`
	GCPProject               = ""         // PROJECT_NAME
	GCPManagedZone           = ""         // DNS_MANAGED_ZONE_NAME
	AddressLines             = []string{} // "Unit 1", "1 One st"
	GoogleRecipients         = []string{} // Company/Individual Name
	RegistrationCity         = ""
	RegistrationCountry      = "" // ISO Country Code i.e AU
	RegistrationPostalCode   = ""
	RegistrationState        = ""
	RegistrationEmail        = ""
	RegistrationPhone        = "" // Format +AreaCode.400000000
	RegistrationFax          = ""
	DnsKind                  = `dns#resourceRecordSet`
	DnsAs                    = []string{} // `216.239.32.21`, `216.239.34.21`, `216.239.36.21`, `216.239.38.21`
	DnsAAAAs                 = []string{} // `2001:4860:4802:32::15`, `2001:4860:4802:34::15`, `2001:4860:4802:36::15`, `2001:4860:4802:38::15`
	RunAdminEndpoint         = `https://us-central1-run.googleapis.com`
	RunAPIVersion            = `serving.knative.dev/v1`
	RunType                  = `Service`
	ProjectNamespace         = `` // namespaces/PROJECT_NAME
	DomainMappingKind        = `DomainMapping`
	DomainMappingNamespace   = `` // PROJECT_NAME
	DomainMappingCertMode    = `AUTOMATIC`
	RunImage                 = `` // us-central1-docker.pkg.dev/PROJET_NAME/cloud-run-source-deploy/PROJET_NAME
)

func main() {

	Dcache = NewDomainCache()
	// Crawl()

	var err error
	fire, err = firebase.NewApp(ctx, nil)
	bug(err)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", Web)
	http.HandleFunc("/intent", buy)
	http.HandleFunc("/purchase", Purchase)
	http.HandleFunc("/cache", GetCache)
	http.HandleFunc("/tlds", GetTlds)
	// http.HandleFunc("/domains", GetDomains)
	http.HandleFunc("/register", TryRegister)
	fmt.Println("waiting")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// Login is done client side with Firebase iOS/Web SDK.
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("uid")) // uuid
	fmt.Println(r.FormValue("did")) // device id
	fmt.Println(r.UserAgent())
}

func Purchase(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 3)
	a := r.URL.Query()
	domains := a["q"]
	for _, a := range domains {
		fmt.Println(a)
	}
}

// HANDLERS

func (d *DomainCache) GetGoogle(word string) []Rjs {
	d.gx.RLock()
	defer d.gx.RUnlock()
	return d.GoogleQueries[word]
}

func (d *DomainCache) SetGoogle(word string, rjs []Rjs) {
	d.gx.Lock()
	defer d.gx.Unlock()
	d.GoogleQueries[word] = rjs
}

// UTILS

func (dc *DomainCache) GetEndings(word string) *Endings {
	dc.x.RLock()
	defer dc.x.RUnlock()
	if w, q := dc.Domains[word]; q {
		return w
	}
	return &Endings{}
}

func (e *Endings) Add(es map[string]interface{}) {
	e.x.Lock()
	defer e.x.Unlock()
	for a, v := range es {
		e.Tlds[a] = v
	}
}

func NewEnding() *Endings {
	return &Endings{
		Tlds: make(map[string]interface{}),
	}
}

func (dc *DomainCache) AddEndings(words string, endings map[string]interface{}) {
	dc.x.Lock()
	defer dc.x.Unlock()
	i, ok := dc.Domains[words]
	if !ok {
		i = NewEnding()
		dc.Domains[words] = i
	}
	i.Add(endings)
}

func (d *DomainCache) ReadCache(word string) []Rjs {
	d.cx.RLock()
	defer d.cx.RUnlock()
	return d.Cache[word]
}

func NewDomainCache() *DomainCache {
	gtls := googletlds()
	gdtls := godaddytlds()
	var k map[string][]Rjs
	// bb, e := ioutil.ReadFile("cache.json")
	// bug(e)
	// e = json.Unmarshal(bb, &k)
	// bug(e)
	// pp.Print(k)
	return &DomainCache{
		Cache: k,
		// Notcache:      t,
		Domains:       make(map[string]*Endings),
		GoogleTlds:    gtls,
		GodaddyTlds:   gdtls,
		Xrates:        Rates(),
		GoogleQueries: make(map[string][]Rjs),
	}
}
