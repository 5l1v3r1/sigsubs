package runner

import (
	"strings"

	"github.com/drsigned/sigsubs/pkg/agent"
	"github.com/drsigned/sigsubs/pkg/sources"
)

// Runner is
type Runner struct {
	Options *Options
	Agent   *agent.Agent
}

// New is
func New(options *Options) *Runner {
	var uses, exclusions []string

	if options.SourcesUse != "" {
		uses = append(uses, strings.Split(options.SourcesUse, ",")...)
	} else {
		uses = append(uses, sources.All...)
	}

	if options.SourcesExclude != "" {
		exclusions = append(exclusions, strings.Split(options.SourcesExclude, ",")...)
	}

	return &Runner{
		Options: options,
		Agent:   agent.New(uses, exclusions),
	}
}

// Run is a
func (runner *Runner) Run() (chan sources.Subdomain, error) {
	// all subdomains
	subdomains := make(chan sources.Subdomain)

	// Create a unique map for filtering out duplicate subdomains
	uniqueMap := make(map[string]sources.Subdomain)
	// Create a map to track source for each subdomain
	sourceMap := make(map[string]map[string]struct{})

	keys := runner.Options.YAMLConfig.GetKeys()
	results := runner.Agent.Run(runner.Options.Domain, &keys)

	go func() {
		defer close(subdomains)

		for result := range results {
			if !strings.HasSuffix(result.Value, "."+runner.Options.Domain) {
				continue
			}

			sub := strings.ToLower(result.Value)

			// remove wildcards (`*`)
			sub = strings.ReplaceAll(sub, "*.", "")

			if _, ok := uniqueMap[sub]; !ok {
				sourceMap[sub] = make(map[string]struct{})
			}

			sourceMap[sub][result.Source] = struct{}{}

			// Check if the sub is a duplicate.
			if _, ok := uniqueMap[sub]; ok {
				continue
			}

			subdomain := sources.Subdomain{Value: sub, Source: result.Source}

			uniqueMap[sub] = subdomain

			subdomains <- subdomain
		}
	}()

	return subdomains, nil
}
