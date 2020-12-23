package sigsubs

import (
	"strings"

	"github.com/drsigned/sigsubs/pkg/sources"
)

// Run is a
func Run(options *Options) (chan sources.Subdomain, error) {
	var uses, exclusions []string

	use := options.UseSources
	exclude := options.ExcludeSources

	// Add sources to use
	if use != "" {
		uses = append(uses, strings.Split(use, ",")...)
	} else {
		uses = append(uses, options.YAMLConfig.Sources...)
	}

	// Add sources to exclude
	if exclude != "" {
		exclusions = append(exclusions, strings.Split(exclude, ",")...)
	}

	passiveAgent := NewAgent(uses, exclusions)

	keys := options.YAMLConfig.GetKeys()
	results := passiveAgent.Enumerate(options.Domain, &keys)

	// Create a unique map for filtering out duplicate subdomains
	uniqueMap := make(map[string]sources.Subdomain)

	// Create a map to track source for each subdomain
	sourceMap := make(map[string]map[string]struct{})

	// all subdomains
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		for result := range results {
			if !strings.HasSuffix(result.Value, "."+options.Domain) {
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
