// Package feedfetch contains types and functions to fetch RSS feeds
package feedfetch

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// RSSFeed is the struct representation of the outer RSS feed tags
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string    `xml:"title"`
		Items []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem is the struct representation of the tags for individual posts in an
// RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

// GetFomURL fetches a single RSS feed of the given URL
func GetFomURL(url string) (RSSFeed, error) {
	rss := RSSFeed{}

	r, err := http.Get(url)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("error fetching -- %v", err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("error reading body -- %v", err)
	}

	xml.Unmarshal(body, &rss)
	return rss, nil
}
