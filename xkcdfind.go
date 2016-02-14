package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alessio/xkcdfind/xkcd"
)

func printResults(results []xkcd.Comic) {
	for _, comic := range results {
		fmt.Printf("%4d %-40s %s\n", comic.Num, comic.Title, comic.Img)
	}
	indexStats := fmt.Sprintf("%d results among %d comics, "+
		"index stats: latest:#%d, missing=%d\n",
		len(results),
		len(xkcd.ComicsIndex.Items),
		xkcd.ComicsIndex.Latest,
		len(xkcd.ComicsIndex.Missing))
	fmt.Fprintf(os.Stderr, indexStats)
}

func main() {
	var (
		indexFilename string
		update        bool
	)
	flag.StringVar(&indexFilename, "index", "", "Index file (default: 'index.json')")
	flag.BoolVar(&update, "update", false, "Force the update of the index")
	flag.Parse()
	if err := xkcd.LoadIndex(indexFilename); err != nil {
		log.Fatal(err)
	}
	if update {
		xkcd.UpdateIndex(indexFilename)
	}
	results := xkcd.RegexSearchComic(flag.Args())
	printResults(results)
}
