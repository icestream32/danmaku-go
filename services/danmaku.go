package services

import (
	"fmt"

	"github.com/gocolly/colly"
)

func GetDanmaku() {

	c := colly.NewCollector(

		colly.AllowedDomains("www.bilibili.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.OnHTML(".body", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {

		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {

		cookies := r.Headers.Get("Set-Cookie")
		fmt.Println(cookies)
	})

	c.OnError(func(r *colly.Response, err error) {

		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	c.Visit("https://www.bilibili.com")
}
