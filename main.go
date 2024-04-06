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
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Tags   string `json:"tags"`
	About  string `json:"about"`
}

// buat detail author
type Detail struct {
	AuthorTitle  string `json:"authortitle"`
	BornDate     string `json:"borndate"`
	BornLocation string `json:"bornlocation"`
	Description  string `json:"description"`
}

// function buat ngambil detail author
func detail() {
	c := colly.NewCollector()

	// details := []Detail{}

	c.OnHTML("", func(h *colly.HTMLElement) {

	})

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

		// fmt.Println(h.ChildAttr("a", "href"))
		q := Quotes{
			Quote:  h.ChildText("span.text"),
			Author: h.ChildText("small.author"),
			Tags:   tagsString,
			About:  h.ChildAttr("a", "href"),
		}
		quotes = append(quotes, q)
	})

	c.OnHTML("span a", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL("href"))
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
