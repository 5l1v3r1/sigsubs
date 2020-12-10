package certspotterv0

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	ID       int      `json:"id"`
	DNSNames []string `json:"dns_names"`
}

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://certspotter.com/api/v0/certs?domain=%s", domain))

		var results []response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results {
			for _, j := range i.DNSNames {
				subdomains <- sources.Subdomain{Source: source.Name(), Value: j}
			}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "certspotterv0"
}
