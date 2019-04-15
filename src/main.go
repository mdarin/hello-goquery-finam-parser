//
// загрузчик данных с сайта Финам
// по мотивам статьи https://habr.com/ru/post/332700/
//
package main

import (
  "fmt"
  "log"

  "net/http"
	"net/url"
	_ "io/ioutil"
	"strings"

	// page parser
  "github.com/PuerkitoBio/goquery"
)



//
// пример использования библиотеки goquery
//
func ExampleScrape() {
  // Request the HTML page.
  res, err := http.Get("http://metalsucks.net")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("a").Text()
    title := s.Find("i").Text()
    fmt.Printf("Review %d: %s - %s\n", i, band, title)
  })
}


func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

func GetAssetsList() {
  // Request the HTML page.
  res, err := http.Get("https://www.finam.ru/quotes/stocks/russia/")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

	// Grabbed document
	//TODO:read doc: https://godoc.org/github.com/PuerkitoBio/goquery 
	// isn't work
	//fmt.Println("document:",doc)

	var CSSPath string

	//CSSPath = "html"

	CSSPath = "html body.i-user_client_no.i-user_client_no div table tbody"


	// Find the table items
  doc.Find(CSSPath).Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
		f := s.Find("td").Text()
		fmt.Println("f:", f)
  })


	fmt.Println()
	fmt.Println("  POST  page 1")
	fmt.Println()

	//We can use POST form to get result, too.
	res, err = http.PostForm("https://www.finam.ru/quotes/stocks/russia/",
		url.Values{"pageNumber": {"1"}})
	if err != nil {
		//panic(err)
    log.Fatal(err)
	}
	defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err = goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

	// Find the table items
  doc.Find(CSSPath).Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
		f := s.Find("td").Text()
		fmt.Println("f:", f)
  })

}

//
// main driver
//
func main() {
	ExampleScrape()

	GetAssetsList()
/*
	// Go contains rich function for grab web contents. net/http is the major library
	// https://dlintw.github.io/gobyexample/public/http-client.html

	// GET request
	resp, err := http.Get("https://dlintw.github.io/gobyexample/public/http-client.html")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", keepLines(string(body), 3))

	//We can use POST form to get result, too.
	resp, err = http.PostForm("http://duckduckgo.com",
		url.Values{"q": {"github"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println("post:\n", keepLines(string(body), 3))
*/
}

