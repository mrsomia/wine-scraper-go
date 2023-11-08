package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type ItemURLs struct {
	tesco     string
	supervalu string
	dunnes    string
}

type ItemLink struct {
	name string
	urls ItemURLs
}

var jd = ItemLink{
    name: "Jack Daniel's", 
    urls: ItemURLs{
      tesco: "https://www.tesco.ie/groceries/en-IE/products/255248604",
      supervalu: "https://shop.supervalu.ie/sm/delivery/rsid/5550/product/jack-daniels-old-no-7-whiskey-70-cl-id-1020340001",
    },
  }

var itemLinks = map[string]ItemLink{
	"JD": jd,
}

func createCollyMap() map[string]*colly.Collector{
	supervaluColly := colly.NewCollector(colly.IgnoreRobotsTxt())
	supervaluColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	supervaluColly.OnHTML(`[data-testid="pdpMainPrice-div-testId"]`, func(e *colly.HTMLElement) {
		fmt.Printf("The Price of JD in Supervalu is: %v", e.Text)
	})

  
  return map[string]*colly.Collector{
    "supervalu": supervaluColly,
  }
}

func main() {
  collys := createCollyMap()
  JDItem := itemLinks["JD"]
  // collys["tesco"].Visit(JDItem.urls.tesco)
  collys["supervalu"].Visit(JDItem.urls.supervalu)
}
