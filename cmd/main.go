package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"github.com/go-rod/rod"
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

type GetFromSite func(url string, wg *sync.WaitGroup)

func cleanPrice(s string) (float64, error) {
	if s == "" {
		return 0.0, fmt.Errorf("No value to convert %v", s)
	}
	i := strings.TrimSuffix(s, " Clubcard Price")
	i = strings.TrimSpace(i)
	i = strings.TrimPrefix(i, "€")
	n, err := strconv.ParseFloat(i, 64)
	if err != nil {
		return 0.0, fmt.Errorf("Unable to convert found value %q to a float. Error: %v", s, err.Error())
	}
	return n, nil
}

func getFromSuperValu(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	el := page.MustElement(`[data-testid="pdpMainPrice-div-testId"]`)
	price, err := cleanPrice(el.MustText())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Price: %v \n", price)
}

func getFromTesco(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	priceEl := page.MustElement(`span[data-auto="price-value"]`)
	offerElem := page.MustElement(`span.offer-text`)
	price, err := cleanPrice(priceEl.MustText())
	if offerElem != nil {
		price, err = cleanPrice(offerElem.MustText())
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Price: %v \n", price)
}

func getFromDunnes(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Fetching from: %v \n", url)
	page := rod.New().MustConnect().MustPage(url)
	page.MustWaitStable()
	el := page.MustElement(`meta[itemprop="price"]`)
	price, err := cleanPrice(*el.MustAttribute("content"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Price: %v \n", price)
}

var scrapers = map[Store]GetFromSite{
	supervalu: getFromSuperValu,
	tesco:     getFromTesco,
	dunnes:    getFromDunnes,
}

func fetchPrices() {
	var wg sync.WaitGroup

	for _, item := range items {
		for _, itemUrl := range item.urls {
			c, ok := scrapers[itemUrl.store]
			if ok {
				wg.Add(1)
				go c(itemUrl.url, &wg)
			}
		}
	}
	wg.Wait()
}

func main() {
  fetchPrices()
}
