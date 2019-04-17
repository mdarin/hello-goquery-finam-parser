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
	"io/ioutil"
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


//TODO: inprove
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

//remove
func GetAssetHistory() {
//	CSSPath := "html body.i-user_logged_no.i-user_client_no div.finam-wrap div.finam-global-container div.content div.layout table.main tbody tr td#content-block.inside-container.content div#issuer-profile div#issuer-profile-container div#issuer-profile-outer div#issuer-profile-inner div#issuer-profile-content div#issuer-profile-export div#issuer-profile-export-form form#chartform.i-form-state"


	CSSPath := "html body.i-user_logged_no.i-user_client_no div.finam-wrap div.finam-global-container div.content div.layout table.main tbody tr td#content-block.inside-container.content div#issuer-profile div#issuer-profile-container div#issuer-profile-outer div#issuer-profile-inner div#issuer-profile-content div#issuer-profile-export div#issuer-profile-export-form"


	//We can use POST form to get result, too.
	res, err := http.PostForm("https://www.finam.ru/profile/moex-akcii/lukoil/export",
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

			s.Find(CSSPath + "form table").Each(func (i int, s1 *goquery.Selection) {
				f := s1.Find("tr").Text()
				fmt.Println("f:",f)

			})
		})
	}
}

/*
Как получить элемент с помощью jQuery?
Для того чтобы понимать как работает селектор Вам все-же необходимы базовые знания CSS, т.к. именно от принципов CSS отталкивается селектор jQuery:

$(“#header”) – получение элемента с id=”header”
$(“h3”) – получить все <h3> элементы
$(“div#content .photo”) – получить все элементы с классом =”photo” которые находятся в элементе div с id=”content”
$(“ul li”) – получить все <li> элементы из списка <ul>
$(“ul li:first”) – получить только первый элемент <li> из списка <ul>
*/
func accessElem() {
	//We can use POST form to get result, too.
	res, err := http.PostForm("https://www.finam.ru/profile/moex-akcii/lukoil/export",
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
		doc.Find("div#issuer-profile-export-form").Each(func (i int, s *goquery.Selection) {
			// For each item found
			id,_ := s.Find("form").Attr("id")
			class,_ := s.Find("form").Attr("class")
			fmt.Println("id:",id)
			fmt.Println("class:",class)
			s.Find("form#"+id+"."+class+" input").Each(func (i int, s1 *goquery.Selection) {
				id1,_ := s1.Attr("id")
				fmt.Println(" id1:",id1)
				name,_ := s.Find("#"+id1).Attr("name")
				fmt.Println("  name:", name)
				value,_ := s.Find("#"+id1).Attr("value")
				fmt.Println("  value:", value)
			})
			s.Find("form#"+id+"."+class+" select").Each(func (i int, s1 *goquery.Selection) {
				id2,_ := s1.Attr("id")
				fmt.Println(" id2:",id2)
				s.Find("#"+id2+" option").Each(func (i int, s2 *goquery.Selection) {
					value,_ := s2.Attr("value")
					fmt.Println("  value:", value)
				})
			})
			s.Find("form#"+id+"."+class+" tbody tr").Each(func (i int, s1 *goquery.Selection) {
				id3,_ := s1.Attr("id")
				if id3 != "" {
					fmt.Println(" id3:",id3)
					text := s.Find("#"+id3+" td").First().Text()
					fmt.Println("text:",text)
				}
			})
		})
	}
}

/*
Для того чтобы написать функцию обращения к серверу «ФИНАМ» (а писать мы будем ее на Python), еще раз рассмотрим параметры GET запроса:
__http://export.finam.ru/POLY_170620_170623.txt?market=1&em=175924&code=POLY&apply=0&df=20&mf=5&yf=2017&from=20.06.2017&dt=23&
mt=5&yt=2017&to=23.06.2017&p=8&f=POLY_170620_170623&e=.txt&cn=POLY&dtf=1&tmf=1&
MSOR=1&mstime=on&mstimever=1&sep=1&sep2=1&datf=1&at=1

POLY_170620_170623 – очевидно, что данная строка представляет параметр code, а также временные характеристики.

.txt – расширение файла; расширение упоминается в параметре e; при написании функции следует помнить об этом нюансе. 

Примем также во внимание содержимое исходного кода страницы типа www.finam.ru/profile/moex-akcii/gazprom/export внутри тэга form (где name=«exportdata»). Характеризуем показатели.

market, em, code – об этих параметрах, упоминал ранее, при обращении к функции их значения будут приниматься из файла.
df, mf, yf, from, dt, mt, yt, to – это параметры времени.
p — период котировок (1 - тики, 2 - 1 мин., 3 - 5 мин., 4 - 10 мин., 5 - 15 мин., 6 -  30 мин., 7 - 1 час, 8 - 1 день,9 - 1 неделя, 10 - 1 месяц)
e – расширение получаемого файла; возможны варианты — .txt либо .csv
dtf — формат даты (1 — ггггммдд, 2 — ггммдд, 3 — ддммгг, 4 — дд/мм/гг, 5 — мм/дд/гг)
tmf — формат времени (1 — ччммсс, 2 — ччмм, 3 — чч: мм: сс, 4 — чч: мм)
MSOR — выдавать время (0 — начала свечи, 1 — окончания свечи)
mstimever — выдавать время (НЕ московское — mstimever=0; московское — mstime='on', mstimever='1')
sep — параметр разделитель полей (1 — запятая (,), 2 — точка (.), 3 — точка с запятой (;), 4 — табуляция (»), 5 — пробел ( ))
sep2 — параметр разделитель разрядов (1 — нет, 2 — точка (.), 3 — запятая (,), 4 — пробел ( ), 5 — кавычка ('))
datf — Перечень получаемых данных (#1 — TICKER, PER, DATE, TIME, OPEN, HIGH, LOW, CLOSE, VOL; #2 — TICKER, PER, DATE, TIME, OPEN, HIGH, LOW, CLOSE; #3 — TICKER, PER, DATE, TIME, CLOSE, VOL; #4 — TICKER, PER, DATE, TIME, CLOSE; #5 — DATE, TIME, OPEN, HIGH, LOW, CLOSE, VOL; #6 — DATE, TIME, LAST, VOL, ID, OPER).
at — добавлять заголовок в файл (0 — нет, 1 — да)
*/
func downloadAssetHistory() {
	// общие параметры
	market := "1"
	em := "8"
	code := "LKOH"
	// неизветный
	apply := "0"
	//параметры времени.
	df := "1"
	mf := "1"
	yf := "2018"
	from := "01.01.2018"
	dt := "17"
	mt := "4"
	yt := "2019"
	to := "17.04.2019"
	//период котировок
	p := "8" // дни
	//расширение получаемого файла
	e := ".csv"
	//формат даты
	dtf := "1"
	//формат времени
	tmf :="1"
	//выдавать время
	MSOR := "0"
	mstimever := "1"
	mstime := "on"
	//параметр разделитель полей
	sep := "3"
	//параметр разделитель разрядов
	sep2 := "2"
	//Перечень получаемых данных
	datf := "1"
	//добавлять заголовок в файл
	at := "1"
	// наименование выходного файла
	f := code + "_" + "180101" + "_" + "190417"

	//Req := "http://export.finam.ru/POLY_170620_170623.txt?market=1&em=175924&code=POLY&apply=0&df=20&mf=5&yf=2017&from=20.06.2017&dt=23&mt=5&yt=2017&to=23.06.2017&p=8&f=POLY_170620_170623&e=.txt&cn=POLY&dtf=1&tmf=1&SOR=1&mstime=on&mstimever=1&sep=1&sep2=1&datf=1&at=1"

	req := "http://export.finam.ru/"+f+e+"?market="+market+"&em="+em+"&code="+code+"&apply="+apply+"&df="+df+"&mf="+mf+"&yf="+yf+"&from="+from+"&dt="+dt+"&mt="+mt+"&yt="+yt+"&to="+to+"&p="+p+"&f="+f+"&e="+e+"&cn="+code+"&dtf="+dtf+"&tmf="+tmf+"&SOR="+MSOR+"&mstime="+mstime+"&mstimever="+mstimever+"&sep="+sep+"&sep2="+sep2+"&datf="+datf+"&at="+at
	fmt.Println("request:",req)


  // Request the HTML page.
  res, err := http.Get(req)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  } else {
		body,_ := ioutil.ReadAll(res.Body)
		fmt.Println("GET:")
		//fmt.Println(keepLines(string(body)))
		fmt.Println(string(body))
	}

/*
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
}


//
// main driver
//
func main() {
//	ExampleScrape()

	//accessElem()
	downloadAssetHistory()
	//GetAssetsList()

	//GetAssetHistory()

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

