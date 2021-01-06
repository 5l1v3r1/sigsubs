package agent

import (
	"sync"

	"github.com/drsigned/sigsubs/pkg/sources"
	"github.com/drsigned/sigsubs/pkg/sources/alienvault"
	"github.com/drsigned/sigsubs/pkg/sources/anubis"
	"github.com/drsigned/sigsubs/pkg/sources/bufferover"
	"github.com/drsigned/sigsubs/pkg/sources/cebaidu"
	"github.com/drsigned/sigsubs/pkg/sources/certspotterv0"
	"github.com/drsigned/sigsubs/pkg/sources/chaos"
	"github.com/drsigned/sigsubs/pkg/sources/crtsh"
	"github.com/drsigned/sigsubs/pkg/sources/github"
	"github.com/drsigned/sigsubs/pkg/sources/hackertarget"
	"github.com/drsigned/sigsubs/pkg/sources/rapiddns"
	"github.com/drsigned/sigsubs/pkg/sources/riddler"
	"github.com/drsigned/sigsubs/pkg/sources/sonar"
	"github.com/drsigned/sigsubs/pkg/sources/sublist3r"
	"github.com/drsigned/sigsubs/pkg/sources/threatcrowd"
	"github.com/drsigned/sigsubs/pkg/sources/threatminer"
	"github.com/drsigned/sigsubs/pkg/sources/urlscan"
	"github.com/drsigned/sigsubs/pkg/sources/wayback"
	"github.com/drsigned/sigsubs/pkg/sources/ximcx"
)

// Agent is a struct for running passive subdomain enumeration
// against a given host. It wraps subsources package and provides
// a layer to build upon.
type Agent struct {
	sources map[string]sources.Source
}

// New creates a new agent for passive subdomain discovery
func New(uses, exclusions []string) *Agent {
	agent := &Agent{
		sources: make(map[string]sources.Source),
	}

	// Add Sources
	for _, source := range uses {
		switch source {
		case "alienvault":
			agent.sources[source] = &alienvault.Source{}
		case "anubis":
			agent.sources[source] = &anubis.Source{}
		case "bufferover":
			agent.sources[source] = &bufferover.Source{}
		case "cebaidu":
			agent.sources[source] = &cebaidu.Source{}
		case "certspotterv0":
			agent.sources[source] = &certspotterv0.Source{}
		case "chaos":
			agent.sources[source] = &chaos.Source{}
		case "crtsh":
			agent.sources[source] = &crtsh.Source{}
		case "github":
			agent.sources[source] = &github.Source{}
		case "hackertarget":
			agent.sources[source] = &hackertarget.Source{}
		case "rapiddns":
			agent.sources[source] = &rapiddns.Source{}
		case "riddler":
			agent.sources[source] = &riddler.Source{}
		case "sonar":
			agent.sources[source] = &sonar.Source{}
		case "sublist3r":
			agent.sources[source] = &sublist3r.Source{}
		case "threatcrowd":
			agent.sources[source] = &threatcrowd.Source{}
		case "threatminer":
			agent.sources[source] = &threatminer.Source{}
		case "urlscan":
			agent.sources[source] = &urlscan.Source{}
		case "wayback":
			agent.sources[source] = &wayback.Source{}
		case "ximcx":
			agent.sources[source] = &ximcx.Source{}
		}
	}

	// Exclude Sources
	for _, source := range exclusions {
		delete(agent.sources, source)
	}

	return agent
}

// Run enumerates all the subdomains for a given domain
func (agent *Agent) Run(domain string, keys *sources.Keys) chan sources.Subdomain {
	results := make(chan sources.Subdomain)

	go func() {
		defer close(results)

		session, _ := sources.NewSession(domain, keys)

		wg := new(sync.WaitGroup)

		// Run each source in parallel on the target domain
		for source, runner := range agent.sources {
			wg.Add(1)

			go func(source string, runner sources.Source) {
				for resp := range runner.Run(domain, session) {
					results <- resp
				}

				wg.Done()
			}(source, runner)
		}

		wg.Wait()
	}()

	return results
}
