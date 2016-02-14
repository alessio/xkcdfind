package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/alessio/xkcdfind/xkcd"
)

const (
	DefaultIndexFilename string = "index.json"
	resultsTemplate             = `{{range .}}{{.Num | printf "%4d"}} {{.Title | printf "%-40s"}} {{.Img}}
{{end}}{{. | len | printf "%4d"}} results
`
)

func printResults(results []xkcd.Comic) {
	report := template.Must(template.New("results").Parse(resultsTemplate))
	if err := report.Execute(os.Stdout, results); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		indexFilename string
		update        bool
	)
	flag.StringVar(&indexFilename, "index", DefaultIndexFilename, "Index file (default: 'index.json')")
	flag.BoolVar(&update, "update", false, "Force the update of the index")
	flag.Parse()
	if err := xkcd.LoadIndex(indexFilename); err != nil {
		log.Fatal(err)
	}
	if update {
		xkcd.UpdateIndex(indexFilename)
	}
	if len(flag.Args()) > 0 {
		results := xkcd.RegexSearchComic(flag.Args())
		printResults(results)
	}
}
