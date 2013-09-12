package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var opts struct {
	Products string `short:"p" long:"products" description:"A YAML file with product data" required:"true"`
	Stores   string `short:"s" long:"stores" description:"A YAML file with store data" required:"true"`
}

var storeList []Store
var productList []Product

type Store struct {
	Name          string
	Domain        string
	XPath         string
	compiledXPath *xpath.Expression
}

type StorePrice struct {
	Store *Store
	Price float64
}

type Product struct {
	Name string
	URLs []string
}

func (p *Product) GetPrices(stores []Store) (store_prices []StorePrice) {
	store_prices = make([]StorePrice, len(p.URLs))

	for i, url := range p.URLs {
		for j, store := range stores {
			if !strings.Contains(url, store.Domain) {
				continue
			}

			resp, err := http.Get(url)
			if err != nil {
				continue
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				continue
			}

			doc, err := gokogiri.ParseHtml(body)

			if err != nil {
				continue
			}

			nxpath := xpath.NewXPath(doc.DocPtr())
			nodes, err := nxpath.Evaluate(doc.DocPtr(), store.compiledXPath)

			if err != nil {
				continue
			}

			if len(nodes) == 0 {
				fmt.Printf("Check XPath correctness (not found) for domain: %s\n", store.Domain)
				continue
			}

			price_raw := xml.NewNode(nodes[0], doc).InnerHtml()
			price_raw = strings.Trim(price_raw, "$ \n\r")
			price, err := strconv.ParseFloat(price_raw, 64)

			if err != nil {
				fmt.Printf("Check XPath correctness (not monetary) for domain: %s\n", store.Domain)
				continue
			}

			store_prices[i] = StorePrice{Store: &stores[j], Price: price}
		}
	}
	return
}

func main() {
	// Parse options
	_, err := flags.Parse(&opts)

	if err != nil {
		fmt.Println("Error: Check options")
		return
	}

	// Open, parse YAML data
	f, _ := os.Open(opts.Products)
	products := make([]byte, 10000)
	count, _ := f.Read(products)
	err = goyaml.Unmarshal(products[:count], &productList)

	if err != nil {
		fmt.Println("Error: Check YAML file of product data")
		return
	}

	f, _ = os.Open(opts.Stores)
	stores := make([]byte, 10000)
	count, _ = f.Read(stores)
	err = goyaml.Unmarshal(stores[:count], &storeList)

	if err != nil {
		fmt.Println("Error: Check YAML file of store data")
		return
	}

	// Compile XPaths
	for i, store := range storeList {
		storeList[i].compiledXPath = xpath.Compile(store.XPath)
	}

	// Loop through products
	for _, product := range productList {
		store_prices := product.GetPrices(storeList)
		fmt.Printf("%s\n", strings.ToUpper(product.Name))

		// Get prices
		for _, store_price := range store_prices {
			fmt.Printf("\t%s: \t$%s\n", store_price.Store.Name, strconv.FormatFloat(store_price.Price, 'f', 2, 32))
		}
	}
}
