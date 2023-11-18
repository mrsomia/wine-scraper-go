package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"strings"
)

type Store int

const (
	tesco Store = iota
	supervalu
	dunnes
)

func (s Store) String() string {
	switch s {
	case tesco:
		return "tesco"
	case supervalu:
		return "supervalu"
	case dunnes:
		return "dunnes"
	default:
		return "unknown"
	}
}

type ItemURL struct {
	store Store
	url   string
}

type Item struct {
	name string
	urls []ItemURL
}

var jd = Item{
	name: "Jack Daniel's",
	urls: []ItemURL{
		{store: tesco, url: "https://www.tesco.ie/groceries/en-IE/products/255248604"},
		{store: supervalu, url: "https://shop.supervalu.ie/sm/delivery/rsid/5550/product/jack-daniels-old-no-7-whiskey-70-cl-id-1020340001"},
		{store: dunnes, url: "https://www.dunnesstoresgrocery.com/sm/delivery/rsid/258/product/jack-daniels-old-no-7-brand-tennessee-sour-mash-whiskey-70cl-id-100672192"},
	},
}

var items = map[string]Item{
	"JD": jd,
}

type GetFromSite func(url string)

func cleanPrice(s string) string {
	i := strings.TrimSuffix(s, " Clubcard Price")
	i = strings.TrimSpace(i)
	i = strings.TrimPrefix(i, "€")
	return i
}

func getFromSuperValu(url string) {
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	el := page.MustElement(`[data-testid="pdpMainPrice-div-testId"]`)
	fmt.Printf("Price: %v \n", cleanPrice(el.MustText()))
}

func getFromTesco(url string) {
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	priceEl := page.MustElement(`span[data-auto="price-value"]`)
	offerElem := page.MustElement(`span.offer-text`)
	if offerElem != nil {
		// offerPrice := strings.Split(offerElem.MustText(), " ")[0]
		fmt.Printf("Price offer: %v \n", cleanPrice(offerElem.MustText()))
		return
	}
	fmt.Printf("Price: %v \n", cleanPrice(priceEl.MustText()))
}

func getFromDunnes(url string) {
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	el := page.MustElement(`meta[itemprop="price"]`)
	fmt.Printf("Price: %v \n", cleanPrice(*el.MustAttribute("content")))
}

var scrapers = map[Store]GetFromSite{
	supervalu: getFromSuperValu,
	tesco:     getFromTesco,
	dunnes:    getFromDunnes,
}

func main() {
	for _, item := range items {
		for _, itemUrl := range item.urls {
			c, ok := scrapers[itemUrl.store]
			if ok {
				c(itemUrl.url)
			}
		}
	}
}
