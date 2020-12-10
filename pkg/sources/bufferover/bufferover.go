package bufferover

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	Meta struct {
		Errors []string `json:"Errors"`
	} `json:"Meta"`
	FDNSA   []string `json:"FDNS_A"`
	RDNS    []string `json:"RDNS"`
	Results []string `json:"Results"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		// Run enumeration on subdomain dataset for historical SONAR datasets
		source.getData(fmt.Sprintf("https://dns.bufferover.run/dns?q=.%s", domain), session, subdomains)
		source.getData(fmt.Sprintf("https://tls.bufferover.run/dns?q=.%s", domain), session, subdomains)
	}()

	return subdomains
}

func (source *Source) getData(sourceURL string, session *sources.Session, subdomains chan sources.Subdomain) {

	res, _ := session.SimpleGet(sourceURL)

	var results response

	if err := json.Unmarshal(res.Body(), &results); err != nil {
		return
	}

	var subs []string

	subs = append(subs, results.FDNSA...)
	subs = append(subs, results.RDNS...)
	subs = append(subs, results.Results...)

	for _, subdomain := range subs {
		for _, value := range session.Extractor.FindAllString(subdomain, -1) {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: value}
		}
	}
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "bufferover"
}
