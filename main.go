package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func feedUrl(event string) string {
	return fmt.Sprintf("https://media.ccc.de/c/%s/podcast/mp4-hq.xml", event)
}

func DownloadFeed(url string) (Rss, error) {
	var result Rss
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = xml.Unmarshal(content, &result)
	return result, err
}

type Rss struct {
	XMLName   xml.Name `xml:"rss"`
	Text      string   `xml:",chardata"`
	Version   string   `xml:"version,attr"`
	Content   string   `xml:"content,attr"`
	Dc        string   `xml:"dc,attr"`
	Trackback string   `xml:"trackback,attr"`
	Itunes    string   `xml:"itunes,attr"`
	Channel   struct {
		Text          string `xml:",chardata"`
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Copyright     string `xml:"copyright"`
		LastBuildDate string `xml:"lastBuildDate"`
		Image         struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			Enclosure   struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Length string `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			PubDate string `xml:"pubDate"`
			Guid    struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Identifier string `xml:"identifier"`
			Date       string `xml:"date"`
			Author     string `xml:"author"`
			Explicit   string `xml:"explicit"`
			Keywords   string `xml:"keywords"`
			Summary    string `xml:"summary"`
			Duration   string `xml:"duration"`
			Subtitle   string `xml:"subtitle"`
		} `xml:"item"`
		Generator string `xml:"generator"`
		Category  struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
		} `xml:"category"`
		Owner struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		Author   string `xml:"author"`
		Explicit string `xml:"explicit"`
		Keywords string `xml:"keywords"`
		Subtitle string `xml:"subtitle"`
		Summary  string `xml:"summary"`
	} `xml:"channel"`
}

func main() {
	url := feedUrl("36c3")

	rss, err := DownloadFeed(url)
	if err != nil {
		panic(err)
	}

	videos := rss.Channel.Item

	for _, v := range videos {
		fmt.Printf("%v\n", v.Title)
	}
}
