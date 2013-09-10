package main

import (
	"os"
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	"launchpad.net/goyaml"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"github.com/moovweb/gokogiri/xml"
	//"launchpad.net/xmlpath"
)

var storeList []Store
var productList []Product

type Store struct {
    Name string
	Domain string
	XPath string
}

type Product struct {
    Name string
	URLs []string
}

func main() {
	f, _ := os.Open("products.yaml")
	products := make([]byte, 10000)
	count, _ := f.Read(products)
	err := goyaml.Unmarshal(products[:count], &productList)

	if err != nil {
		panic(err)
	}

	f, _ = os.Open("stores.yaml")
	stores := make([]byte, 10000)
	count, _ = f.Read(stores)
	err = goyaml.Unmarshal(stores[:count], &storeList)

	if err != nil {
		panic(err)
	}

	for _, product := range productList {
		getPrices(product, storeList)
	}

	fmt.Printf("hi")
}

func getPrices(product Product, stores []Store) (price float32, error string) {
	for _, url := range product.URLs {
		for _, store := range stores {
			if(!strings.Contains(url, store.Domain)){
				continue
			}

			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			doc, err := gokogiri.ParseHtml(body)
			exp := xpath.Compile(store.XPath)
			nxpath := xpath.NewXPath(doc.DocPtr())
			nodes, err := nxpath.Evaluate(doc.DocPtr(), exp)
			if(len(nodes) > 0){
				price := xml.NewNode(nodes[0], doc).InnerHtml()
				price = strings.Trim(price, " \n\r")
				fmt.Printf("%s", price)
			}
		}
	}
	return
}
