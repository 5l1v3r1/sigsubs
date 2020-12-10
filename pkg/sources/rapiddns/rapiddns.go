package rapiddns

import (
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://rapiddns.io/subdomain/%s?full=1", domain))

		for _, subdomain := range session.Extractor.FindAllString(string(res.Body()), -1) {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: subdomain}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "rapiddns"
}
