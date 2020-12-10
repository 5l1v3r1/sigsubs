package alienvault

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	Detail     string `json:"detail"`
	Error      string `json:"error"`
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain))

		var results response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		if results.Error != "" {
			return
		}

		for _, j := range results.PassiveDNS {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: j.Hostname}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "alienvault"
}
