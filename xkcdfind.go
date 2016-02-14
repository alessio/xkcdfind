package main

import (
	"flag"
	"fmt"
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

func mustLoadIndex(filename string) *xkcd.Index {
	index, err := xkcd.LoadIndex(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	return index
}

func main() {
	var (
		filename string
		index    *xkcd.Index
		update   bool
	)
	flag.StringVar(&filename, "index", DefaultIndexFilename, "Index file")
	flag.BoolVar(&update, "update", false, "Force the update of the index")
	flag.Parse()

	if !update && flag.NArg() == 0 {
		index = mustLoadIndex(filename)
		fmt.Fprintf(os.Stderr, "%s\n", index)
		os.Exit(0)
	}
	if update {
		index, err := xkcd.LoadIndex(filename)
		if err != nil {
			if os.IsNotExist(err) {
				index = new(xkcd.Index)
			} else {
				fmt.Fprintf(os.Stderr, "%s is corrupted -- %s\n", filename, err)
				os.Exit(2)
			}
		}
		if err := index.UpdateIndex(filename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		fmt.Fprintf(os.Stderr, "%s\n", index)
	}
	if flag.NArg() > 0 {
		index = mustLoadIndex(filename)
		results := index.RegexSearchComic(flag.Args())
		printResults(results)
	}
}
