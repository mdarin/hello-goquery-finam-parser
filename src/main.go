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
	var CSSPath string

/*
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



	 Find the table items
  doc.Find(CSSPath).Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
		f := s.Find("td").Text()
		fmt.Println("f:", f)
  })

*/

	CSSPath = "html body.i-user_client_no.i-user_client_no div table tbody"

	stop := false
	for i := 1; i < 100 && !stop; i++ {
		fmt.Println()
		fmt.Println("  POST  page ",i)
		fmt.Println()

		// convert to string
		page := fmt.Sprintf("%d", i)

		//We can use POST form to get result, too.
		res, err := http.PostForm("https://www.finam.ru/quotes/stocks/russia/",
			url.Values{"pageNumber": {page}})
		if err != nil {
			//panic(err)
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			stop = true
		} else {
			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Find the table items
			doc.Find(CSSPath).Each(func (i int, s *goquery.Selection) {
				// For each item found, get the band and title
				s.Find("tr").Each(func (i int, s1 *goquery.Selection) {
					href,_ := s1.Find("a").Attr("href")
					title,_ := s1.Find("a").Attr("title")
					fmt.Println("href:", href)
					fmt.Println("title:", title)

					s1.Find("td").Each( func (i int, s2 *goquery.Selection) {
						f := s2.Find("span").Text()
						fmt.Println("f:", f)
					})
				})
			})
		}
	}

}


func GetAssetHistory() {
//	CSSPath := "html body.i-user_logged_no.i-user_client_no div.finam-wrap div.finam-global-container div.content div.layout table.main tbody tr td#content-block.inside-container.content div#issuer-profile div#issuer-profile-container div#issuer-profile-outer div#issuer-profile-inner div#issuer-profile-content div#issuer-profile-export div#issuer-profile-export-form form#chartform.i-form-state"


	CSSPath := "html body.i-user_logged_no.i-user_client_no div.finam-wrap div.finam-global-container div.content div.layout table.main tbody tr td#content-block.inside-container.content div#issuer-profile div#issuer-profile-container div#issuer-profile-outer div#issuer-profile-inner div#issuer-profile-content div#issuer-profile-export div#issuer-profile-export-form"


	//We can use POST form to get result, too.
	res, err := http.PostForm("https://www.finam.ru/profile/moex-akcii/sberbank/export",
		nil)//url.Values{"pageNumber": {page}})
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	} else {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Find the table items
		doc.Find(CSSPath).Each(func (i int, s *goquery.Selection) {
			// For each item found
			id,_ := s.Find(CSSPath + " form").Attr("id")
			class,_ := s.Find(CSSPath + " form").Attr("class")
			name,_ := s.Find(CSSPath + " form").Attr("name")
			action,_ := s.Find(CSSPath + " form").Attr("action")
			method,_ := s.Find(CSSPath + " form").Attr("method")
			fmt.Println("id:",id)
			fmt.Println("class:",class)
			fmt.Println("name:",name)
			fmt.Println("action:",action)
			fmt.Println("method:",method)
			s.Find(CSSPath + " form input").Each(func (i int, s1 *goquery.Selection) {
				id,_ := s1.Attr("id")
				//typ,_ := s1.Attr("type")
				name,_ := s1.Attr("name")
				value,_ := s1.Attr("value")
				fmt.Println(" *****")
				fmt.Println(" id:",id)
				//fmt.Println(" type:",typ)
				fmt.Println(" name:",name)
				fmt.Println(" value:",value)
			})
		})
	}
}


//
// main driver
//
func main() {
//	ExampleScrape()

	//GetAssetsList()

	GetAssetHistory()
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

