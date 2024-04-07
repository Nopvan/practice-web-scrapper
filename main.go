package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// buat quotes
type Quotes struct {
	Quote   string `json:"quote"`
	Author  string `json:"author"`
	Tags    string `json:"tags"`
	About   string `json:"about"`
	Details Detail `json:"details"`
}

// buat detail author
type Detail struct {
	AuthorTitle  string `json:"authortitle"`
	BornDate     string `json:"borndate"`
	BornLocation string `json:"bornlocation"`
	Description  string `json:"description"`
}

// function buat ngambil detail author
func details(authorURL string) Detail {
	// authorURL = ("https://quotes.toscrape.com/author/Albert-Einstein/")
	c := colly.NewCollector()

	var d Detail
	// details := []Detail{}

	c.OnHTML("div.author-details", func(h *colly.HTMLElement) {
		d = Detail{
			AuthorTitle:  h.ChildText("h3.author-title"),
			BornDate:     h.ChildText("span.author-born-date"),
			BornLocation: h.ChildText("span.author-born-location"),
			Description:  h.ChildText("div.author-description"),
		}
		// details = append(details, d)
	})

	// Visit authorURL untuk mengambil detail penulis
	c.Visit(authorURL)

	return d

}

func main() {
	c := colly.NewCollector()

	quotes := []Quotes{}

	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnHTML("div.quote", func(h *colly.HTMLElement) {

		//buat spasi di bagian tags dengan cara for each , iterate class tag, terus tag.value + ' ', terus bikin array, tinggal push ke array, abis itu array nya tinggal di masukin ke json
		//shoutout to ojan and chat gpt :3
		var tags []string
		h.ForEach("a.tag", func(i int, h *colly.HTMLElement) {
			tags = append(tags, h.Text)
		})

		tagsString := strings.Join(tags, " ")

		// buat ambil URL penulis dari "About" link
		authorUrl := h.ChildAttr("a", "href")

		// buat ambil detail penulis menggunakan fungsi details
		authorDetail := details("https://quotes.toscrape.com" + authorUrl)

		// fmt.Println(h.ChildAttr("a", "href"))
		q := Quotes{
			Quote:   h.ChildText("span.text"),
			Author:  h.ChildText("small.author"),
			Tags:    tagsString,
			About:   authorUrl,
			Details: authorDetail,
		}
		quotes = append(quotes, q)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://quotes.toscrape.com/page/1/")

	// fmt.Println(quotes)
	data, err := json.Marshal(quotes)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("quotes.json", data, 0644)
}
