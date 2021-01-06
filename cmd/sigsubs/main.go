package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/drsigned/sigsubs/pkg/runner"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	sourcesList bool
	noColor     bool
	silent      bool
	verbosity   int
}

var (
	co options
	so runner.Options
	au aurora.Aurora
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                 _         
 ___(_) __ _ ___ _   _| |__  ___ 
/ __| |/ _`+"`"+` / __| | | | '_ \/ __|
\__ \ | (_| \__ \ |_| | |_) \__ \
|___/_|\__, |___/\__,_|_.__/|___/ V1.3.0
       |___/
`).Bold())
}

func init() {
	flag.StringVar(&so.Domain, "d", "", "")
	flag.StringVar(&so.SourcesExclude, "sE", "", "")
	flag.BoolVar(&co.sourcesList, "sL", false, "")
	flag.BoolVar(&co.noColor, "nC", false, "")
	flag.BoolVar(&co.silent, "silent", false, "")
	flag.StringVar(&so.SourcesUse, "sU", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  sigsubs [OPTIONS]\n"

		h += "\nOPTIONS:\n"
		h += "  -d          domain to find subdomains for\n"
		h += "  -sE         comma(,) separated list of sources to exclude\n"
		h += "  -sL         list all the sources available\n"
		h += "  -nC         no color mode\n"
		h += "  -silent     silent mode: output subdomains only\n"
		h += "  -sU         comma(,) separated list of sources to use\n\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	options, err := runner.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	if !co.silent {
		banner()
	}

	if co.sourcesList {
		fmt.Println("[", au.BrightBlue("INF"), "] current list of the available", au.Underline(strconv.Itoa(len(options.YAMLConfig.Sources))+" sources").Bold())
		fmt.Println("[", au.BrightBlue("INF"), "] sources marked with an * needs key or token")
		fmt.Println("")

		keys := options.YAMLConfig.GetKeys()
		needsKey := make(map[string]interface{})
		keysElem := reflect.ValueOf(&keys).Elem()

		for i := 0; i < keysElem.NumField(); i++ {
			needsKey[strings.ToLower(keysElem.Type().Field(i).Name)] = keysElem.Field(i).Interface()
		}

		for _, source := range options.YAMLConfig.Sources {
			if _, ok := needsKey[source]; ok {
				fmt.Println(">", source, "*")
			} else {
				fmt.Println(">", source)
			}
		}

		os.Exit(0)
	}

	if !co.silent {
		fmt.Println("[", au.BrightBlue("INF"), "]", au.Underline(so.Domain).Bold(), "subdomain enumeration.")
		fmt.Println("")
	}

	runner := runner.New(options)

	subdomains, err := runner.Run()
	if err != nil {
		log.Fatalln(err)
	}

	for n := range subdomains {
		if co.silent {
			fmt.Println(n.Value)
		} else {
			fmt.Println(fmt.Sprintf("[%s] %s", au.BrightBlue(n.Source), n.Value))
		}
	}
}
