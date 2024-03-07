package feedfetch

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

func GetFomURL(url string) (rssFeed, error) {
	rss := rssFeed{}

	r, err := http.Get(url)
	if err != nil {
		return rssFeed{}, fmt.Errorf("error fetching -- %v", err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return rssFeed{}, fmt.Errorf("error reading body -- %v", err)
	}

	xml.Unmarshal(body, &rss)
	return rss, nil
}
