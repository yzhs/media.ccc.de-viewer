package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

type Item struct {
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
		Item      []Item `xml:"item"`
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

type Video struct {
	Event     string `json:"event"`
	Id        string `json:"id"`
	Title     string `json:"string"`
	ShortLink string `json:"short_link"`
	Link      string `json:"link"`
}

func GetVideoList(event string, url string) []Video {
	rss, err := DownloadFeed(url)
	if err != nil {
		log.Fatal(err)
	}

	rawVideos := rss.Channel.Item

	videos := make([]Video, len(rawVideos))

	for i, v := range rawVideos {
		title := v.Title
		if v.Subtitle != "" {
			title += " — " + v.Subtitle
		}

		tags := strings.Split(v.Keywords, ", ")

		video := Video{event, tags[1], title, v.Link, v.Enclosure.URL}
		videos[i] = video
	}

	return videos
}

func PrintVideo(video Video) {
	fmt.Println(video.Id, video.Title)
	fmt.Println(video.ShortLink)
	fmt.Println(video.Link)
	fmt.Println()
}

func PrintJson(video Video) {
	v, err := json.Marshal(video)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", v)
}

func ListVideos(event string) {
	url := feedUrl(event)

	videos := GetVideoList(event, url)

	for _, video := range videos {
		PrintJson(video)
	}
}

func Usage() {
	fmt.Printf("Usage: %s list <event>\n", os.Args[0])
}

func main() {
	if len(os.Args) < 3 {
		Usage()
		return
	}

	command := os.Args[1]
	if command != "list" {
		Usage()
		return
	}

	event := os.Args[2]
	ListVideos(event)
}
