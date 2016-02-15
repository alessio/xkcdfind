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
		log.Fatalf("Couldn't print results: %s", err)
	}
}

func mustLoadIndex(filename string) *xkcd.Index {
	index, err := xkcd.LoadIndex(filename)
	if err != nil {
		log.Fatalf("Couldn't load the index: %s", err)
	}
	return index
}

func main() {
	var (
		filename string
		index    *xkcd.Index
		update   bool
	)
	log.SetFlags(0)
	flag.StringVar(&filename, "index", DefaultIndexFilename, "Index file")
	flag.BoolVar(&update, "update", false, "Force the update of the index")
	flag.Parse()

	if !update && flag.NArg() == 0 {
		index = mustLoadIndex(filename)
		log.Printf("%s", index)
		os.Exit(0)
	}
	if update {
		index, err := xkcd.LoadIndex(filename)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("The index file '%s' is corrupted: %s", filename, err)
			}
			index = new(xkcd.Index)
		}
		if err := index.UpdateIndex(filename); err != nil {
			log.Fatalf("Couldn't update the index: %s", err)
		}
		log.Printf("%s", index)
	}
	if flag.NArg() > 0 {
		index = mustLoadIndex(filename)
		results := index.RegexSearchComic(flag.Args())
		printResults(results)
	}
}
