package github

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/drsigned/sigsubs/pkg/sources"
	"github.com/tomnomnom/linkheader"
	"github.com/valyala/fasthttp"
)

// Source is the passive sources agent
type Source struct{}

type textMatch struct {
	Fragment string `json:"fragment"`
}

type item struct {
	Name        string      `json:"name"`
	HTMLURL     string      `json:"html_url"`
	TextMatches []textMatch `json:"text_matches"`
}

type response struct {
	TotalCount int    `json:"total_count"`
	Items      []item `json:"items"`
}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		if len(session.Keys.GitHub) == 0 {
			return
		}

		tokens := NewTokenManager(session.Keys.GitHub)

		searchURL := fmt.Sprintf("https://api.github.com/search/code?per_page=100&q=%s&sort=created&order=asc", domain)
		source.Enumerate(searchURL, domainRegexp(domain), tokens, session, subdomains)
	}()

	return subdomains
}

// Enumerate is a
func (source *Source) Enumerate(searchURL string, domainRegexp *regexp.Regexp, tokens *Tokens, session *sources.Session, subdomains chan sources.Subdomain) {
	token := tokens.Get()

	if token.RetryAfter > 0 {
		if len(tokens.pool) == 1 {
			time.Sleep(time.Duration(token.RetryAfter) * time.Second)
		} else {
			token = tokens.Get()
		}
	}

	// Initial request to GitHub search
	res, err := session.Request(
		fasthttp.MethodGet,
		searchURL,
		"",
		map[string]string{
			"Accept":        "application/vnd.github.v3.text-match+json",
			"Authorization": "token " + token.Hash,
		},
		nil,
	)
	isForbidden := res != nil && res.StatusCode() == fasthttp.StatusForbidden
	if err != nil && !isForbidden {
		return
	}

	// Retry enumerarion after Retry-After seconds on rate limit abuse detected
	ratelimitRemaining, _ := strconv.ParseInt(string(res.Header.Peek("X-Ratelimit-Remaining")), 10, 64)
	if isForbidden && ratelimitRemaining == 0 {
		retryAfterSeconds, _ := strconv.ParseInt(string(res.Header.Peek("Retry-After")), 10, 64)
		tokens.setCurrentTokenExceeded(retryAfterSeconds)

		source.Enumerate(searchURL, domainRegexp, tokens, session, subdomains)
	}

	var results response

	// Marshall json response
	if err := json.Unmarshal(res.Body(), &results); err != nil {
		return
	}

	err = proccesItems(results.Items, domainRegexp, source.Name(), session, subdomains)
	if err != nil {
		return
	}

	// Links header, first, next, last...
	linksHeader := linkheader.Parse(string(res.Header.Peek("Link")))
	// Process the next link recursively
	for _, link := range linksHeader {
		if link.Rel == "next" {
			nextURL, err := url.QueryUnescape(link.URL)
			if err != nil {
				return
			}
			source.Enumerate(nextURL, domainRegexp, tokens, session, subdomains)
		}
	}
}

// proccesItems procceses github response items
func proccesItems(items []item, domainRegexp *regexp.Regexp, name string, session *sources.Session, results chan sources.Subdomain) error {
	for _, item := range items {
		// find subdomains in code
		res, _ := session.SimpleGet(rawURL(item.HTMLURL))

		if res.StatusCode() == fasthttp.StatusOK {
			scanner := bufio.NewScanner(bytes.NewReader(res.Body()))
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				for _, subdomain := range domainRegexp.FindAllString(normalizeContent(line), -1) {
					results <- sources.Subdomain{Source: name, Value: subdomain}
				}
			}
		}

		// find subdomains in text matches
		for _, textMatch := range item.TextMatches {
			for _, subdomain := range domainRegexp.FindAllString(normalizeContent(textMatch.Fragment), -1) {
				results <- sources.Subdomain{Source: name, Value: subdomain}
			}
		}
	}
	return nil
}

// Normalize content before matching, query unescape, remove tabs and new line chars
func normalizeContent(content string) string {
	normalizedContent, _ := url.QueryUnescape(content)
	normalizedContent = strings.ReplaceAll(normalizedContent, "\\t", "")
	normalizedContent = strings.ReplaceAll(normalizedContent, "\\n", "")
	return normalizedContent
}

// Raw URL to get the files code and match for subdomains
func rawURL(htmlURL string) string {
	domain := strings.ReplaceAll(htmlURL, "https://github.com/", "https://raw.githubusercontent.com/")
	return strings.ReplaceAll(domain, "/blob/", "/")
}

// DomainRegexp regular expression to match subdomains in github files code
func domainRegexp(domain string) *regexp.Regexp {
	rdomain := strings.ReplaceAll(domain, ".", "\\.")
	return regexp.MustCompile("(\\w+[.])*" + rdomain)
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "github"
}
