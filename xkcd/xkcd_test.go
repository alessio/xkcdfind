package xkcd

import (
	"fmt"
	"testing"
)

func TestFetchComic(t *testing.T) {
	comicID := 5

	comic, err := FetchComic(comicID)
	fmt.Println(comic.Num)
	fmt.Println(err)
	// Output:
	// 5
	// <nil>
}
