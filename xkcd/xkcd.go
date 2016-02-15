package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	XkcdURL            string = "http://xkcd.com/"
	RemoteJSONFilename string = "info.0.json"
)

type Comic struct {
	Num        int
	SafeTitle  string `json:"safe_title"`
	Alt        string
	Img        string
	Title      string
	Transcript string
}

type Index struct {
	Items   map[string]Comic
	Latest  int
	Missing []int
}

func (ind *Index) String() string {
	return fmt.Sprintf("comics:%d  latest:#%d  missing:%d",
		len(ind.Items), ind.Latest, len(ind.Missing))
}

func LoadIndex(filename string) (*Index, error) {
	var ind Index

	fp, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(fp).Decode(&ind)
	if err != nil {
		return nil, err
	}
	return &ind, nil
}

func (ind *Index) UpdateIndex(filename string) error {
	latestRemoteComic, err := FetchComic(0) // Fetch latest
	if err != nil {
		return fmt.Errorf(
			"couldn't retrieve remote's latest comic -- %s", err)
	}
	if ind.Latest == 0 {
		ind.Items = make(map[string]Comic)
	}

	for i := ind.Latest + 1; i <= latestRemoteComic.Num; i++ {
		if comic, err := FetchComic(i); err != nil {
			fmt.Fprintf(os.Stderr,
				"couldn't retrieve comic -- %s\n", err)
			ind.Missing = append(ind.Missing, i)
		} else {
			ind.Items[strconv.Itoa(i)] = *comic
			ind.Latest = i
		}
	}

	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer fp.Close()
	if err != nil {
		return fmt.Errorf("couldn't open '%s' -- %s", filename, err)
	}

	return json.NewEncoder(fp).Encode(ind)
}

func (ind *Index) RegexSearchComic(terms []string) []Comic {
	var (
		results []Comic
		rs      []*regexp.Regexp
	)

	for _, expr := range terms {
		if r, err := regexp.Compile(expr); err == nil {
			rs = append(rs, r)
		} else {
			fmt.Fprintf(os.Stderr, "Invalid regex: %s\n", expr)
		}
	}
	for _, comic := range ind.Items {
		for _, r := range rs {
			if r.FindStringIndex(comic.Alt) != nil ||
				r.FindStringIndex(comic.Title) != nil ||
				r.FindStringIndex(comic.SafeTitle) != nil ||
				r.FindStringIndex(comic.Transcript) != nil {
				results = append(results, comic)
				break
			}
		}
	}
	return results
}

func FetchComic(comicID int) (*Comic, error) {
	var (
		comic Comic
		url   string
	)

	if comicID == 0 {
		url = strings.Join([]string{XkcdURL, RemoteJSONFilename}, "")
	} else {
		url = strings.Join([]string{XkcdURL, strconv.Itoa(comicID), "/", RemoteJSONFilename}, "")
	}

	fmt.Fprintf(os.Stderr, "Fetching remote index: %s\n", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"couldn't fetch comic '%d' -- %d", comicID, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return nil, err
	}

	return &comic, nil
}
