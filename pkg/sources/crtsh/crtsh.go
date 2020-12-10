package crtsh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	ID        int    `json:"id"`
	NameValue string `json:"name_value"`
}

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain))

		var results []response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results {
			x := strings.Split(i.NameValue, "\n")

			for _, j := range x {
				subdomains <- sources.Subdomain{Source: source.Name(), Value: j}
			}

		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "crtsh"
}
