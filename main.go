package main

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly/v2"
  "github.com/go-rod/rod"
)

type ItemURL struct {
	store string
	url   string
}

type ItemLink struct {
	name string
	urls []ItemURL
}

var jd = ItemLink{
	name: "Jack Daniel's",
	urls: []ItemURL{
		{store: "tesco", url: "https://www.tesco.ie/groceries/en-IE/search?query=jack%20daniel"},
		{store: "supervalu", url: "https://shop.supervalu.ie/sm/delivery/rsid/5550/product/jack-daniels-old-no-7-whiskey-70-cl-id-1020340001"},
    {store: "dunnes", url: "https://www.dunnesstoresgrocery.com/sm/delivery/rsid/258/product/jack-daniels-old-no-7-brand-tennessee-sour-mash-whiskey-70cl-id-100672192"},
	},
}

var itemLinks = map[string]ItemLink{
	"JD": jd,
}

const USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"

type GetFromSite func(url string) 

func getFromSuperValu(url string) {
	supervaluColly := colly.NewCollector(colly.IgnoreRobotsTxt())
  supervaluColly.UserAgent = USER_AGENT
	supervaluColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	supervaluColly.OnHTML(`[data-testid="pdpMainPrice-div-testId"]`, func(e *colly.HTMLElement) {
    fmt.Printf("The Price of JD in Supervalu is: %v\n", e.Text[1:])
	})

  supervaluColly.Visit(url)
}

func getFromTesco(url string) {
	tescoColly := colly.NewCollector(colly.IgnoreRobotsTxt())
	tescoColly.UserAgent = USER_AGENT


	tescoColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

  tescoColly.OnHTML(`div[data-auto="product-tile"]`, func(e *colly.HTMLElement) {
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
  
  tescoColly.Visit(url)
}

func getFromDunnes(url string) {
  page := rod.New().MustConnect().MustPage(url)
  page.MustWaitStable()
  el := page.MustElement(`meta[itemprop="price"]`)
  fmt.Println(*el.MustAttribute("content"))
}


var scrapers = map[string]GetFromSite{
		"supervalu": getFromSuperValu,
		"tesco":     getFromTesco,
		"dunnes":    getFromDunnes,
	}

func main() {
  for _, itemLink := range itemLinks {
    for _, itemUrl := range itemLink.urls {
      c, ok := scrapers[itemUrl.store]
      if ok {
        c(itemUrl.url)
      }
    }
  }
}
