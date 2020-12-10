package chaos

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
	"github.com/valyala/fasthttp"
)

type response struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Count      int      `json:"count"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		if session.Keys.Chaos == "" {
			return
		}

		res, _ := session.Request(
			fasthttp.MethodGet,
			fmt.Sprintf("https://dns.projectdiscovery.io/dns/%s/subdomains", domain),
			"",
			map[string]string{"Authorization": session.Keys.Chaos},
			nil,
		)

		var results response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results.Subdomains {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: fmt.Sprintf("%s.%s", i, results.Domain)}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "chaos"
}
