package main

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly/v2"
)


func main() {
	c := colly.NewCollector(colly.IgnoreRobotsTxt())
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.360"

  URL := "https://www.tesco.ie/groceries/en-IE/search?query=jack%20daniels"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

  c.OnHTML(`div[data-auto="product-tile"]`, func(e *colly.HTMLElement) {
    tile := e.DOM
    h3 := tile.Find("h3").Text()
    if h3 != "Jack Daniel's Old No. 7 Tennessee Whiskey 70 cL" {
      return
    }

    price := tile.Find(".beans-price__container  p").First().Text()
    offerPrice := tile.Find("span.offer-text").First().Text()
    if offerPrice != "" {
      price = fmt.Sprintf("%v", strings.Split(offerPrice, " ")[0])
    }
    fmt.Printf("Price of %v at Tesco is: %v \n", h3, price[1:])
	})

  c.Visit(URL)
}
