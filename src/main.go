//
// загрузчик данных с сайта Финам
// по мотивам статьи https://habr.com/ru/post/332700/
//
// Work like a slave; command like a king; create like a god.
// Original in Romanian: 
// Muncește ca un sclav, poruncește ca un rege, creează ca un zeu.
// Constantin Brâncuși
package main

import (
	// common purpose
	"fmt"
	"log"

	// statndard
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"bufio" // to scan and tokenize buffered input data from an io.Reader source
	"strconv"
	"regexp"
	"errors" // for errors.New()
	"os" // for operations with dirs
	"time" // for sleep
	"sort" // for sorging
	//NOTE: The path package should only be used for paths separated by forward slashes, 
	//      such as the paths in URLs. This package does not deal with Windows paths with 
	//      drive letters or backslashes; to manipulate operating system paths, use the path/filepath package.
	"path"

	// Package filepath implements utility routines for manipulating filename paths in a way
	// compatible with the target operating system-defined file paths.
	"path/filepath"

	// page parser
	"github.com/PuerkitoBio/goquery"
	// toUtf8 conversin
	// https://stackoverflow.com/questions/6927611/go-language-how-to-convert-ansi-text-to-utf8
	"golang.org/x/text/encoding/charmap"
)

//TODO: обязательно добавить таймауты на заросы!!
// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779


// по регэкспам https://shapeshed.com/golang-regexp/


const(

)

var(
	counter = 0
)

// Win1251toUtf8 and Win1251fromUtf8 conversin
// https://stackoverflow.com/questions/6927611/go-language-how-to-convert-ansi-text-to-utf8
func DecodeWindows1251(enc []byte) string {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

func EncodeWindows1251(inp string) []byte {
	enc := charmap.Windows1251.NewEncoder()
	out, _ := enc.String(inp)
	// For converting from a string to a byte slice, string -> []byte:
	// https://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
	// []byte(str)
	return []byte(out)
}



func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}


// получить ссылки и наименования активов(здесь акций РФР)
func getAssetsList(dir string) {
	stop := false
	for i := 1; i < 100 && !stop; i++ {

		// convert to string
		page := fmt.Sprintf("%d", i)

		//We can use POST form to get result, too.
		res, err := http.PostForm("https://www.finam.ru/quotes/stocks/russia/",
			url.Values{"pageNumber": {page}})
		if err != nil {
			//panic(err)
			log.Fatal(err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			stop = true
			fmt.Println("  DONE ")
		} else {

			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			CSSPath := "html body.i-user_client_no.i-user_client_no div table tbody"

			// Find the table items
			doc.Find(CSSPath).Each(func (i int, tbody *goquery.Selection) {
				// For each item found, get the band and title
				tbody.Find("tr").Each(func (i int, tr *goquery.Selection) {
					// получить ссылку и наименование актива
					href,_ := tr.Find("a").Attr("href")
					title,_ := tr.Find("a").Attr("title")
					//convert to utf-8
					title = DecodeWindows1251([]byte(title))
					fmt.Println()
					fmt.Println("page:["+page+"]")
					fmt.Println("tr href:", href)
					fmt.Println("tr title:", title)
					// цифры текущего состояния актива
					//tr.Find("td").Each( func (i int, td *goquery.Selection) {
					//	spanValue := td.Find("span").Text()
					//	fmt.Println("span value:", spanValue)
					//})
					//получить параметры актива перейдя по ссылке и выбрав требуемые значения полей
					getAssetParams(href,title,dir)
					// выдержка перед следующем запросом
					time.Sleep(1000 * 2 * time.Millisecond)
				})
			})
		}
	}
}

/*
Как получить элемент с помощью jQuery?
Для того чтобы понимать как работает селектор Вам все-же необходимы базовые знания CSS, т.к. 
именно от принципов CSS отталкивается селектор jQuery:

$(“#header”) – получение элемента с id=”header”
$(“h3”) – получить все <h3> элементы
$(“div#content .photo”) – получить все элементы с классом =”photo” которые находятся в элементе div с id=”content”
$(“ul li”) – получить все <li> элементы из списка <ul>
$(“ul li:first”) – получить только первый элемент <li> из списка <ul>
*/
//func accessElem() {
//}

// получить параметры по каждому интсрументу
func getAssetParams(href,title,dir string) {

	fmt.Println(" ::href:",href)
	fmt.Println(" ::title:",title)
	fmt.Println(" ::dir:",dir)

	counter++
	fmt.Println(" ::count:",counter)

	// сформировать запросо
	req := path.Join("www.finam.ru",href,"export")
	req = "https://" + req
	fmt.Println(" ::req:",req)

	//We can use POST form to get result, too.
	//res, err := http.PostForm("https://www.finam.ru/profile/moex-akcii/lukoil/export",nil)
	res, err := http.PostForm(req,nil)
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	} else {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// параметры для формирования запроса
		params := make(map[string]string)
		// Find the table items
		doc.Find("div#issuer-profile-export-form").Each(func (i int, form *goquery.Selection) {
			// For each item found
			formId,_ := form.Find("form").Attr("id")
			formClass,_ := form.Find("form").Attr("class")
			fmt.Println("form id:",formId)
			fmt.Println("form class:",formClass)
			// значения параметров запроса из полей формы ввода
			form.Find("form#"+formId+"."+formClass+" input").Each(func (i int, input *goquery.Selection) {
				inputId,_ := input.Attr("id")
				//fmt.Println(" input id:",inputId)
				inputName,_ := form.Find("#"+inputId).Attr("name")
				//fmt.Println(" input name:", inputName)
				inputValue,_ := form.Find("#"+inputId).Attr("value")
				//fmt.Println(" input value:", inputValue)
				// записать параметр
				params[inputName] = inputValue
			})

			// парамтеры выпадающих списков
			/*
			form.Find("form#"+formId+"."+formClass+" select").Each(func (i int, sel *goquery.Selection) {
				selectId,_ := sel.Attr("id")
				fmt.Println(" select:",selectId)
				form.Find("#"+selectId+" option").Each(func (i int, option *goquery.Selection) {
					value,_ := option.Attr("value")
					fmt.Println(" select value:", value)
				})
			})
			*/
			// названия строк формы
			/*
			form.Find("form#"+formId+"."+formClass+" tbody tr").Each(func (i int, tr *goquery.Selection) {
				trId,_ := tr.Attr("id")
				// костыль, в выводе есть пустой тэг tr, это условие его фильтрует
				if trId != "" {
					fmt.Println(" tr id:",trId)
					text := form.Find("#"+trId+" td").First().Text()
					fmt.Println(" tr text:",DecodeWindows1251([]byte(text)))
				}
			})
			*/
		})

		// задать натройки по умолчанию либо жёстко заданные
		//params["market"] = "1"
		//params["em"] = "8"
		//params["code"] = "LKOH"
		// неизветный
		params["apply"] = "0"
		//параметры времени.
		params["df"] = "1"
		params["mf"] = "1"
		params["yf"] = "2009"//"2018"
		params["from"] = strings.Join([]string{params["df"],params["mf"],params["yf"]},".")//"01.01.2018"
		//params["dt"] = "17"
		//params["mt"] = "4"
		//params["yt"] = "2019"
		//params["to"] = "17.04.2019"
		//период котировок
		params["p"] = "8" // дни
		//расширение получаемого файла
		params["e"] = ".csv"
		//формат даты
		params["dtf"] = "1"
		//формат времени
		params["tmf"] ="1"
		//выдавать время
		params["MSOR"] = "0" // MSOR0 ??
		params["mstimever"] = "1" // MSOR1 ??
		params["mstime"] = "on" // пустое значение приходит
		//параметр разделитель полей
		params["sep"] = "3"
		//параметр разделитель разрядов
		params["sep2"] = "2"
		//Перечень получаемых данных
		params["datf"] = "1"
		//добавлять заголовок в файл
		params["at"] = "1"
		// каталог для расположения файла котировок
		params["dir"] = dir

		// загрузить истоирю инструмента с указанными параметрами	
		downloadAssetHistory(params)
	}
}

/*
Для того чтобы написать функцию обращения к серверу «ФИНАМ», еще раз рассмотрим параметры GET запроса:

http://export.finam.ru/POLY_170620_170623.txt?
market=1
&em=175924
&code=POLY
&apply=0
&df=20
&mf=5
&yf=2017
&from=20.06.2017
&dt=23
&mt=5
&yt=2017
&to=23.06.2017
&p=8
&f=POLY_170620_170623
&e=.txt
&cn=POLY
&dtf=1
&tmf=1
&MSOR=1
&mstime=on
&mstimever=1
&sep=1
&sep2=1
&datf=1
&at=1

POLY_170620_170623 – очевидно, что данная строка представляет параметр code, а также временные характеристики.
.txt – расширение файла; расширение упоминается в параметре e; 

NOTE: при написании функции следует помнить об этом нюансе. 

Примем также во внимание содержимое исходного кода страницы типа www.finam.ru/profile/moex-akcii/gazprom/export внутри тэга form (где name=«exportdata»). Характеризуем показатели.

Среди всего перечня хотелось бы акцентировать внимание на параметрах em, market, code. 

em — параметр следует понимать как индекс, своеобразную метку бумаги (инструмента).
		Если мы хотим скачивать не один инструмент, а массив данных по нескольким бумагам 
		k(инструментам) мы должны знать em каждого из них. 
market — говорит о том, где вращается данная бумага (инструмент) – на каком рынке? 
		Маркетов много: МосБиржа топ***, МосБиржа пифы***, МосБиржа облигации***, Расписки и т.д. 
code – это символьная переменная по инструменту. 
df, mf, yf, from, dt, mt, yt, to – это параметры времени.
p — период котировок (
		1 - тики, 
		2 - 1 мин.,
		3 - 5 мин.,
		4 - 10 мин.,
		5 - 15 мин.,
		6 -  30 мин.,
		7 - 1 час,
		8 - 1 день,
		9 - 1 неделя,
		10 - 1 месяц)
e – расширение получаемого файла( 
		.txt
		.csv)
dtf — формат даты (
		1 — ггггммдд,
		2 — ггммдд,
		3 — ддммгг,
		4 — дд/мм/гг,
		5 — мм/дд/гг)
tmf — формат времени (
		1 — ччммсс,
		2 — ччмм,
		3 — чч: мм: сс,
		4 — чч: мм)
MSOR — выдавать время (
		0 — начала свечи,
		1 — окончания свечи)
mstimever — выдавать время (
		НЕ московское — mstimever=0; 
		московское — mstime='on', 
		mstimever='1')
sep — параметр разделитель полей (
		1 — запятая (,),
		2 — точка (.),
		3 — точка с запятой (;),
		4 — табуляция (»),
		5 — пробел ( ))
sep2 — параметр разделитель разрядов (
		1 — нет,
		2 — точка (.),
		3 — запятая (,),
		4 — пробел ( ),
		5 — кавычка ('))
datf — Перечень получаемых данных (FIXME: венрно ли это здесь?
		1 — TICKER, PER, DATE, TIME, OPEN, HIGH, LOW, CLOSE, VOL; 
		2 — TICKER, PER, DATE, TIME, OPEN, HIGH, LOW, CLOSE; 
		3 — TICKER, PER, DATE, TIME, CLOSE, VOL; 
		4 — TICKER, PER, DATE, TIME, CLOSE; 
		5 — DATE, TIME, OPEN, HIGH, LOW, CLOSE, VOL; 
		6 — DATE, TIME, LAST, VOL, ID, OPER).
at — добавлять заголовок в файл (
		0 — нет, 
		1 — да)
FIXME: надо проверить!
fsp — Заполнять периоды без сделок(
		0 — нет, 
		1 — да)

*/
func downloadAssetHistory(params map[string]string) {

	// map traversal
	for k,v := range params {
		fmt.Println("parameter ",k,":",v)
	}


	// общие параметры
	market := params["market"]//"1"
	em := params["em"]//"8"
	code := params["code"]//"LKOH"
	// неизветный
	apply := params["apply"]//"0"
	// параметры времени.
	df := params["df"]//"1"
	mf := params["mf"]//"1"
	yf := params["yf"]//"2018"
	from := params["from"]//"01.01.2018"
	dt := params["dt"]//"17"
	mt := params["mt"]//"4"
	yt := params["yt"]//"2019"
	to := params["to"]//"17.04.2019"
	// период котировок
	p := params["p"]//"8" // дни
	// расширение получаемого файла
	e := params["e"]//".csv"
	// формат даты
	dtf := params["dtf"]//"1"
	// формат времени
	tmf :=params["tmf"]//"1"
	// выдавать время
	MSOR := params["MSOR"]//"0"
	mstimever := params["mstimever"]//"1"
	mstime := params["mstime"]//"on"
	// параметр разделитель полей
	sep := params["sep"]//"3"
	// параметр разделитель разрядов
	sep2 := params["sep2"]//"2"
	// Перечень получаемых данных
	datf := params["datf"]//"1"
	// добавлять заголовок в файл
	at := params["at"]//"1"
	// Заполнять периоды без сделок
	fsp := params["fsp"]//"1"
	// наименование выходного файла
	// https://golang.org/pkg/regexp/#pkg-examples
	re := regexp.MustCompile(`(?P<day>[0-9]+)[.](?P<month>[0-9]+)[.][0-9]{2,2}(?P<year>[0-9]{2,2})`)
	fromPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	toPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	fromName := re.ReplaceAllString(from, fromPart)
	toName := re.ReplaceAllString(to, toPart)
	f := code + "_" + fromName + "_" + toName

	//fmt.Println("from:",re.MatchString(from))
	//fmt.Println("to:",re.MatchString(to))
	//fmt.Println("fromPart:",re.MatchString(from))
	//fmt.Println("toPart:",re.MatchString(to))
	//fmt.Println("fromName:",fromName)
	//fmt.Println("toName:",toName)

	//TODO: можно собрать и из хэша, но надо аккуратно см.пробник mapEx(), там есть нюансы возможно порядок следования аргументов имеет значение
	// запрос истории иснтрумента с указанными параметрами
	req := "http://export.finam.ru/"+f+e+"?market="+market+"&em="+em+"&code="+code+"&apply="+apply+"&df="+df+"&mf="+mf+"&yf="+yf+"&from="+from+"&dt="+dt+"&mt="+mt+"&yt="+yt+"&to="+to+"&p="+p+"&f="+f+"&e="+e+"&cn="+code+"&dtf="+dtf+"&tmf="+tmf+"&SOR="+MSOR+"&mstime="+mstime+"&mstimever="+mstimever+"&sep="+sep+"&sep2="+sep2+"&datf="+datf+"&at="+at+"&fsp="+fsp
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
		// создать файл для размещения котировок актива
		assetFname := filepath.Join(params["dir"],f+e)
		asset,err := os.Create(assetFname)
		if err != nil {
			log.Fatal(err)
		}
		defer asset.Close()
		asset.WriteString(string(body))
	}

}




// подготовить структуру каталогов для размещения скачиваемых котировок
func prepare() string {
	// fetching pwd
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("working dir is:",workdir)
	dirName := "stocks"
	dirPath := filepath.Join(workdir,dirName)
	fmt.Println("dir:",dirPath)
	//TODO: remove all or what to do with existing directory?
	// if data directory is not exists then create it at once
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Println("directory \""+dirPath+"\" does not exist")
		fmt.Println("Creating \""+dirPath+"\"...")
		//ModePerm FileMode = 0777 // Unix permission bits
		var perm os.FileMode = 0777
		os.Mkdir(dirPath, perm)
	}
	return dirPath
}



// ----------------
type SummaryTableRecord struct{
	//ticker string // тикер акции
	date string // дата фиксации цены
	price float64 // цена закрытия
}
// constructor
func NewSTR(date string, price float64) (*SummaryTableRecord,error) {
	rec := new(SummaryTableRecord)
	rec.date = date
	rec.price = price
	return rec,nil
}
type SummaryTable map[string][]*SummaryTableRecord //string
// интерфейс
// присоединить с головы 
func (s *SummaryTable) InsertBefore(ticker,date string, price float64) {
	rec,_ := NewSTR(date, price)
	(*s)[ticker] = append([]*SummaryTableRecord{rec},(*s)[ticker]...)
}
// присоединить с хвоста
func (s *SummaryTable) InsertAfter(ticker,date string, price float64) {
	rec,_ := NewSTR(date, price)
	(*s)[ticker] = append((*s)[ticker], rec)
}

type TickersAndLens struct {
	length int
	ticker string
}
// This type implements sort.Interface for []TicketsAndLens based on that length values.
type ByLen []TickersAndLens

func (a ByLen) Len() int           { return len(a) }
func (a ByLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// asc
//func (a ByLen) Less(i, j int) bool { return a[i].length < a[j].length }
// decs
func (a ByLen) Less(i, j int) bool { return a[i].length > a[j].length }

// ------------------


// заполнить сводную таблицу котировок данными из загруженных файлов
func fillin(rootDir string, summaryTable *SummaryTable) error {
	count := 0
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			//fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			fmt.Printf("Is a dir: %+v %q\n", info.Name(), path)
			//return filepath.SkipDir
		} else {
			fmt.Println()
			fmt.Printf("visited file: %q  %d\n", path, count)
			// выбрать из файла котировок(quotations) данные тикер, дата, цена закрытия 
			// и поместить их в буфер для последующего формирвоания сводной таблицы
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Unable to open file:", path, err)
				return err
			}
			defer file.Close()

			fmt.Println("TIKER:",filepath.Base(path))
			once := false

			// размер таблицы до изменений 
			before := len( (*summaryTable) )

			fmt.Println("Summary len before:",before)
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			// skip head
			// пропускать заголовок :)
			scanner.Scan()
			for scanner.Scan() {
				//FIXME:DEBUG показать запись
				//fmt.Println("$$ ",scanner.Text())

				//FIXME:DEBUG показать запись
				//fmt.Println("*",fields)

				// получить поля
				fields := strings.Split(scanner.Text(), ";")
				// tiker:field[0], date:field[2], close:field[7]
				// тикер актива
				ticker := fields[0]
				// показать тикер однократно
				if !once {
					fmt.Printf("tiker: %s\n",ticker)
					once = true
				}
				// дата
				date := fields[2]
				// преобразовать в float64
				price,_ := toFloat64(fields[7])
				//fmt.Printf("tiker: %s date: %q  price: %q[%.4f]\n",fields[0],fields[2],fields[7],floatNum)
				// https://stackoverflow.com/questions/12677934/create-a-golang-map-of-lists
				//summaryTable[ticker] = append(summaryTable[ticker],date)
				(*summaryTable).InsertAfter(ticker,date,price)
			}
			// размер таблицы после изменений 
			after := len( (*summaryTable) )
			fmt.Println("Summary len after:",after)
			// если размер действительно изменился, то увеличить счётчик
			if before < after {
				count++
			}
		} // eof if-else

		return nil
	})

	return err
}


// отсортировать наполненную таблицу по длине историй
func getSortedByLen(summaryTable *SummaryTable, byLen *ByLen) {
	// сформировать вспомогательный массив для сортировки по длине истории,
	// чтобы получить значение максимальной длины
	// и попутно упорядоченный список длин историй инструментов
	i := 0
	// map traversal
	for k,v := range (*summaryTable) {
		//fmt.Printf("ticker: %s  len:%d sumlen:%d\n",k,len(v),len(summaryTable))
		(*byLen)[i].ticker = k
		(*byLen)[i].length = len(v)
		i++
	}
	// получить наиболее длинную историю за одно отсортировать истории
	sort.Sort(byLen)

	//FIXME:DEBUG
	// отобразить полученные результат
	for i := 0; i < len( (*byLen) ); i++ {
		fmt.Printf("ticker: %s  len:%d   sumlen:%d\n",(*byLen)[i].ticker,(*byLen)[i].length,(*byLen).Len())
		//fmt.Printf("[%d]: ticker: %s  len:%d\n",i,byLen[i].ticker,byLen[i].length)
	}
} // eof func


// выровнять таблицу, дополняя нулями попуски для выравнивания таблицы
func align(summaryTable *SummaryTable, byLen *ByLen, longestTicker string) {
	// выровнять размеры, добивая пропуски нулями для всех инструментов
	// тут надо соотносить по датам значеиния цены и забивать нулями всё пробелы и пропуски
	//
	//	              / 0 если date[i,j] не существует или dete[i,j] младше date[i,0] модельной даты
	//	price[i,j] = <
	//	              \ price[i,j] если date[i,j] == date[i,0] - дате наиболее длинной истории
	// выбрать модельную дату
	// сравинить модельную дату и дату инструмента
	// если дата инструмента младше модельной либо её вообще нет то установить значение цены на эту дату 0(нуль)
	// если дата инструмента равна модельной дате, то установить значение цены на эту дату как цену инструмента считанную из файла котировок инструмента
	// для каждой строки выражающей дату(день) отобразить значение цены каждого инструмента на эту дату
	for i := 0; i < (*byLen)[0].length; i++ {
		// получить дату текущую(дату итерации) из наиболее длинной истории
		strModelDate := (*summaryTable)[longestTicker][i].date
		modelDate,_ := toDate(strModelDate)

		// для каждого актива
		for _,asset := range (*byLen) {
			// получить дату интсрумента
			strAssetDate := (*summaryTable)[asset.ticker][i].date
			assetDate,_ := toDate(strAssetDate)

			// если дата i-го элемента истории младше или старше модельного,
			// то внести запись в историю на эту дату с нулевой ценой в голову или хвост истории
			// иначе пропустить(оставить) актуальную цену для текущего актива
			if modelDate.Before(assetDate) {
				price := float64(0)
				(*summaryTable).InsertBefore(asset.ticker,strModelDate,price)
			} else if assetDate.Before(modelDate) {
				price := float64(0)
				(*summaryTable).InsertAfter(asset.ticker,strModelDate,price)
			} else if i+1 == len( (*summaryTable)[asset.ticker]) {
				// на всякий случай чтоб за границы не выйти, а то бывало
				price := float64(0)
				(*summaryTable).InsertAfter(asset.ticker,strModelDate,price)
			}
		}
	}
} // eof func


// формирование сводной таблицы
func build(summaryTable *SummaryTable, byLen *ByLen, longestTicker string) {
	// сформирвоать заголовок
	// поле дата
	fmt.Printf("DATE")
	// получить значение тикера для каждого инструмента
	for _,asset := range (*byLen) {
		fmt.Printf("%s%s",";",asset.ticker)
	}
	fmt.Println()

	// для каждой строки выражающей дату(день) отобразить значение цены каждого инструмента на эту дату
	for i := 0; i < (*byLen)[0].length; i++ {
		// получить дату текущую(дату итерации) из наиболее длинной истории
		date := (*summaryTable)[longestTicker][i].date
		//TODO:(вроде готово)вообще надо нормировку дат произвести чтобы вдруг пустых не было пропуски 
		// все надо нулями забить в модельной истории перед использованием
		// отображить дату
		sep := ""
		fmt.Printf("%s%s",sep,date)
		sep = ";"
		// получить значение цены на дату для каждого инструмента
		for _,asset := range (*byLen) {
			//fmt.Printf("%s%s",sep,asset.ticker)
			// тут вообще не должно быть условий, пробегаться по всей строке и всё
			// воводть цену данного актива на текущую(согласно итерации) дату наболее длинной истории
			// таблица должна быть уже отформатирована, т.е. все инструменты должны 
			// быть соотнесены по дате и добиты нулями
			price := (*summaryTable)[asset.ticker][i].price
			fmt.Printf("%s%f",sep,price)
		}
		fmt.Println()
	}
} // eof func


// привести скачанные данные к виду пригодному для загрузки в решающее устройство
// здесь надо сформировать на выходе сводную таблицу активов(Assets), в которой все пробелы дополнены нулями
func transform(rootDir string) {

	// сводная таблица активов
	// строка таблицы - это значения каждого актива на указанную дату либо 0 если значение на дату отсутствует
	summaryTable := make(SummaryTable)

	// заполнить сводную таблицу котировок данными из загруженных файлов
	if err := fillin(rootDir, &summaryTable); err != nil {
		//fmt.Printf("error walking the path %q: %v\n", tmpDir, err)
		fmt.Println("transform error:",err)
		return
	}

	// временный массив для сортировки
	byLen := make(ByLen,len(summaryTable))

	// отсортировать наполненную таблицу по длине историй во вспомогательный массив
	getSortedByLen(&summaryTable,&byLen)

	// получить наибольшую длину и тикер с самой длинной историей
	longestTicker := byLen[0].ticker
	longetsLen := byLen[0].length

	// показать наиболее длинную иторию
	fmt.Printf("MAX ticker: %s  len:%d\n",longestTicker,longetsLen)

	// 
	//TODO:25.04.2019
	// пока это двухпроходная схема обработки
	// на первом проходе таблица форматируется таким образом чтобы в строке содеражались 
	// значения цен инструментов на дату или нули если истории нет
	// на втором проходе формируется итоговая сводная таблица содержащая в себе только
	// требуемые данные пригодные для загрузки в матричную бибилиотеку
	//

	// выровнять таблицу, дополняя нулями попуски для выравнивания таблицы
	//align(&summaryTable, &byLen, longestTicker)

	// сформировать сводную таблицу
	//build(&summaryTable, &byLen, longestTicker)
} // eof func


// convert Finam quotations' prices value to Golang float64
func toFloat64(s string) (float64,error) {
	re := regexp.MustCompile(`^(?P<integer>.+?)[.](?P<fractional>[0-9]+)$`)
	//fmt.Println("string:",s)
	if re.MatchString(s) {
		//fmt.Println("string:",s)
		// split fractional part by colon ':'
		temp := re.ReplaceAllString(s,"${integer}:${fractional}")
		//fmt.Println(temp)
		// remove all dots '.'
		re = regexp.MustCompile(`[.]`)
		temp = re.ReplaceAllString(temp,"")
		//fmt.Println(temp)
		// replace colon by dot (to well formed float)
		re = regexp.MustCompile(`[:]`)
		temp = re.ReplaceAllString(temp,".")
		//fmt.Println(temp)
		// convert to float64
		floatNum, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			return 0.0,err
			//fmt.Println("error:",err)
		}
		//fmt.Printf("float: %f\n", floatNum)
		return floatNum,nil
	} else {
		//fmt.Println("not match",s)
		return 0.0,errors.New("not match")
	}
}


/*
Алгоритм

надо будет разобраться со структурой чтобы по каталогам это всё аккуратно заложить


TODO: для каждого рынка[Акции,Облигации?] 
	для каждого иструмента 
		-получить страницу для загрузки данных истории
		-на странице загрузки итории получить параметры требуемы для загрузки
		-сформировать запрос для загрузки данных и загрузить данные истории
	сформировать сводную таблицу котировок всех активов по ценам закрытия соотнесённых по времени(дате)
*/

//
// main driver
//
func main() {

	markets := []string{
		"https://www.finam.ru/quotes/stocks/russia/", //- Акции российкий фондовый рынок
		"https://www.finam.ru/quotes/indices/", // - Индексы
		"https://www.finam.ru/quotes/bonds/", // - Облигации
	}

	for _,m := range markets {
		fmt.Println("market:",m)
		list := strings.Split(m,"/")
		fmt.Println("list:",list[len(list)-3:])
	}
	dataDir := prepare()
	//getAssetsList(dataDir)
	transform(dataDir)


//
// пробники
//

//	cmpDateEx()

	//strToDateEx()
	//appNprepEx()

	//toFloatEx()

	//getAssetParams()
	//downloadAssetHistory()

	//downloadEx()
	//dirEx()
	//mapEx()

}


//
// пробники
//


// go parse date and time
// тут шаблоны описаны вроде норм
// http://demin.ws/blog/russian/2012/04/27/date-and-time-formatting-in-go/
func strToDateEx() {
	// преобразумеа строка содержащая дату
	str := []string{
		"20130806",
		"20140212",
		"20150402",
		"20160216",
		"20180306",
		"20091201",
	}
	// эти волшебшые числа в шаблоне это пи*да, именно так должно быть!"
	layout := "20060102" // формат разбираемой даты
	// полуичить дату 
	for _,s := range str {
		date, err := time.Parse(layout,s)
		if err != nil {
			panic(err)
		}
		fmt.Println("date:",date)
	}
}

// преобразовать строку даты в формате Финан в дату Go
func toDate(s string) (time.Time,error) {
	// эти волшебшые числа в шаблоне это пи*да, именно так должно быть!"
	layout := "20060102" // формат разбираемой даты
	date, err := time.Parse(layout,s)
	if err != nil {
		return date,err
	}
	return date,nil
}





// Go – append/prepend item into slice
func appNprepEx() {
	// https://codingair.wordpress.com/2014/07/18/go-appendprepend-item-into-slice/
	data := []string{"A", "B", "C", "D"}
	// append "F"
  data = append(data, "F")
  fmt.Println("append:",data)
  // [A B C D F]
	// prepend "G"
  data = append([]string{"G"},data...)
  fmt.Println("prepend:",data)
  // [G A B C D F]


	/*
	тамже в коментах
	A better prepend, as it generates less garbage:
	
	data = append(data, “”)
	copy(data[1:], data)
	data[0] = “Prepend Item”
	*/
}


func toFloatEx() {
	str := []string{
		"0.2250000",
		"0.2220000",
		"0.2200000",
		"0.2120000",
		"0.2150000",
		"0.2240000",
		"0.2200000",
		"0.2200000",
		"0.2210000",
		"714.0000000",
		"701.0000000",
		"682.0000000",
		"750.0000000",
		"700.0000000",
		"722.0000000",
		"680.0000000",
		"717.0000000",
		"749.0000000",
		"903.0000000",
		"894.0000000",
		"1.088.0000000",
		"1.040.0000000",
		"1.074.0000000",
		"1.051.0000000",
		"1.132.0000000",
		"1.186.0000000",
		"1.129.0000000",
		"1.102.0000000",
		"1.099.0000000",
		"1.078.0000000",
		"1.089.0000000",
		"1.035.0000000",
		"1.034.0000000",
		"1.032.0000000",
		"990.0000000",
		"1.002.0000000",
		"1.005.0000000",
		"a;dsf",
	}

	/*
	from := "01.01.2019"
	re := regexp.MustCompile(`(?P<integer>[0-9]+)[.](?P<fractional>[0-9]+)`)
	fromPart := fmt.Sprintf("${%s}${%s}", re.SubexpNames()[1], re.SubexpNames()[2])
	fromName := re.ReplaceAllString(from, fromPart)

	fmt.Println("from:",re.MatchString(from))
	fmt.Println("fromPart:",re.MatchString(from))
	fmt.Println("fromName:",fromName)
	*/

	// integer and fractional part of the malformed float price value
	for	_,s := range str {
		fmt.Println("string:",s)
		re := regexp.MustCompile(`^(?P<integer>.+?)[.](?P<fractional>[0-9]+)$`)
		// split fractional part by colon ':'
		temp := re.ReplaceAllString(s,"${integer}:${fractional}")
		fmt.Println(temp)
		// remove all dots '.'
		re = regexp.MustCompile(`[.]`)
		temp = re.ReplaceAllString(temp,"")
		fmt.Println(temp)
		// replace colon by dot (to well formed float)
		re = regexp.MustCompile(`[:]`)
		temp = re.ReplaceAllString(temp,".")
		fmt.Println(temp)
		// convert to float64
		floatNum, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			//return 0.0,nil
			fmt.Println("error:",err)
		}
		fmt.Printf("float: %f\n", floatNum)
	}
}

func mapEx() {
	params := make(map[string]string)
	params["market"] = "1"
  params["em"] = "8"
  params["code"] = "LKOH"
  // неизветный
  params["apply"] = "0"
  //параметры времени.
  params["df"] = "1"
  params["mf"] = "1"
  params["yf"] = "2018"
  params["from"] = "01.01.2018"
  params["dt"] = "17"
  params["mt"] = "4"
  params["yt"] = "2019"
  params["to"] = "17.04.2019"
  //период котировок
  params["p"] = "8" // дни
  //расширение получаемого файла
  params["e"] = ".csv"
  //формат даты
  params["dtf"] = "1"
  //формат времени
  params["tmf"] ="1"
  //выдавать время
  params["SOR"] = "0"
  params["mstimever"] = "1"
  params["mstime"] = "on"
  //параметр разделитель полей
  params["sep"] = "3"
  //параметр разделитель разрядов
  params["sep2"] = "2"
  //Перечень получаемых данных
  params["datf"] = "1"
  //добавлять заголовок в файл
  params["at"] = "1"
	re := regexp.MustCompile(`(?P<day>[0-9]+)[.](?P<month>[0-9]+)[.][0-9]{2,2}(?P<year>[0-9]{2,2})`)
	fromPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	toPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	fromName := re.ReplaceAllString(params["from"], fromPart)
	toName := re.ReplaceAllString(params["to"], toPart)
	params["f"] = params["code"] + "_" + fromName + "_" + toName
	t := params["f"] + params["e"]

	// сформировать запрос истории иснтрумента с указанными параметрами
	req := "http://export.finam.ru/" + t
	sep := "?"
	// Map traversal
	// The for...range loop statement can be used to walk the content of a map value
	for k,v := range params {
		// технически можно, но надо внимательно следить, не все пары совпадают
		// например, &SOR=0 в настоящем запроса
		// а у тут &MSOR=0
		// в целом выгрузка происходит, 
		// но это опасненько, надо отлаживать
		req += sep + k + "=" + v
		sep = "&"
	}
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
}

func downloadEx() {

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
	// https://golang.org/pkg/regexp/#pkg-examples
	re := regexp.MustCompile(`(?P<day>[0-9]+)[.](?P<month>[0-9]+)[.][0-9]{2,2}(?P<year>[0-9]{2,2})`)
	fromPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	toPart := fmt.Sprintf("${%s}${%s}${%s}", re.SubexpNames()[3], re.SubexpNames()[2], re.SubexpNames()[1])
	fromName := re.ReplaceAllString(from, fromPart)
	toName := re.ReplaceAllString(to, toPart)
	f := code + "_" + fromName + "_" + toName
	//fmt.Println("from:",re.MatchString(from))
	//fmt.Println("to:",re.MatchString(to))
	//fmt.Println("fromPart:",re.MatchString(from))
	//fmt.Println("toPart:",re.MatchString(to))
	//fmt.Println("fromName:",fromName)
	//fmt.Println("toName:",toName)

	// http://export.finam.ru/POLY_170620_170623.txt?market=1&em=175924&code=POLY&apply=0&df=20&mf=5&yf=2017&from=20.06.2017&dt=23&mt=5&yt=2017&to=23.06.2017&p=8&f=POLY_170620_170623&e=.txt&cn=POLY&dtf=1&tmf=1&SOR=1&mstime=on&mstimever=1&sep=1&sep2=1&datf=1&at=1

	// запрос истории иснтрумента с указанными параметрами
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

}


func dirEx() {
	// fetching pwd
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("working dir is:",workdir)
	// existance and filepath joining
	filename := filepath.Join(workdir,"a-nonexistent-file.txt")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("file \""+filename+"\" does not exist")
	}

	fmt.Println("args[0]:",os.Args[0])
	fmt.Println("dir:",filepath.Dir(os.Args[0]))
	fmt.Println("base:",filepath.Base(os.Args[0]))


	dirName := "stocks"
	dirPath := filepath.Join(workdir,dirName)
	fmt.Println("dir:",dirPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Println("directory \""+dirPath+"\" does not exist")
		fmt.Println("Creating \""+dirPath+"\"...")
		//ModePerm FileMode = 0777 // Unix permission bits
		perm := 0777
		os.Mkdir(dirPath, os.FileMode(perm))
	}

	//os.RemoveAll(dirPath) // удаляет вместе с корневым каталогом

	// обход от указанного корня
	// https://golang.org/pkg/path/filepath/#Walk
	//subDirToSkip := "skip"
	tmpDir := "dir/to/walk/skip"
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() { // && info.Name() == subDirToSkip {
			//fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			fmt.Printf("Is a dir: %+v \n", info.Name())
			//return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", tmpDir, err)
		return
	}

}


func cmpDateEx() {
	model := []string{
"20090202",
"20090203",
"20090204",
"20090205",
"20090206",
"20090209",
"20090210",
"20090211",
"20090212",
"20090213",
"20090216",
"20090217",
"20090218",
"20090219",
"20090220",
"20090224",
"20090225",
"20090226",
"20090227",
"20090302",
"20090303",
"20090304",
"20090305",
"20090306",
"20090310",
"20090311",
"20090312",
"20090313",
"20090316",
"20090317",
"20090318",
"20090319",
"20090320",
"20090323",
"20090324",
"20090325",
"20090326",
"20090327",
"20090330",
"20090331",
"20090401",
"20090402",
"20090403",
"20090406",
"20090407",
"20090408",
"20090409",
"20090410",
"20090413",
"20090414",
"20090415",
"20090416",
"20090417",
"20090420",
"20090421",
"20090422",
"20090423",
"20090424",
"20090427",
"20090428",
"20090429",
"20090430",
"20090504",
"20090505",
"20090506",
"20090507",
"20090508",
"20090512",
"20090513",
"20090514",
"20090515",
"20090518",
"20090519",
"20090520",
"20090521",
"20090522",
"20090525",
"20090526",
"20090527",
"20090528",
"20090529",
"20090601",
"20090602",
"20090603",
"20090604",
"20090605",
"20090608",
"20090609",
"20090610",
"20090611",
"20090615",
"20090616",
"20090617",
"20090618",
"20090619",
"20090622",
"20090623",
"20090624",
"20090625",
"20090626",
"20090629",
"20090630",
"20090701",
"20090702",
"20090703",
"20090706",
"20090707",
"20090708",
"20090709",
"20090710",
"20090713",
"20090714",
"20090715",
"20090716",
"20090717",
"20090720",
"20090721",
"20090722",
"20090723",
"20090724",
"20090727",
"20090728",
"20090729",
"20090730",
"20090731",
"20090803",
"20090804",
"20090805",
"20090806",
"20090807",
"20090810",
"20090811",
"20090812",
"20090813",
"20090814",
"20090817",
"20090818",
"20090819",
"20090820",
"20090821",
"20090824",
"20090825",
"20090826",
"20090827",
"20090828",
"20090831",
"20090901",
"20090902",
"20090903",
"20090904",
"20090907",
"20090908",
"20090909",
"20090910",
"20090911",
"20090914",
"20090915",
"20090916",
"20090917",
"20090918",
"20090921",
"20090922",
"20090923",
"20090924",
"20090925",
"20090928",
"20090929",
"20090930",
"20091001",
"20091002",
"20091005",
"20091006",
"20091007",
"20091008",
"20091009",
"20091012",
"20091013",
"20091014",
"20091015",
"20091016",
"20091019",
"20091020",
"20091021",
"20091022",
"20091023",
"20091026",
"20091027",
"20091028",
"20091029",
"20091030",
"20091102",
"20091103",
"20091105",
"20091106",
"20091109",
"20091110",
"20091111",
"20091112",
"20091113",
"20091116",
"20091117",
"20091118",
"20091119",
"20091120",
"20091123",
"20091124",
"20091125",
"20091126",
"20091127",
"20091130",
"20091201",
"20091202",
"20091203",
"20091204",
"20091207",
"20091208",
"20091209",
"20091210",
"20091211",
"20091214",
"20091215",
"20091216",
"20091217",
"20091218",
"20091221",
"20091222",
"20091223",
"20091224",
"20091225",
"20091228",
"20091229",
"20091230",
"20091231",
"20100111",
"20100112",
"20100113",
"20100114",
"20100115",
"20100118",
"20100119",
"20100120",
"20100121",
"20100122",
"20100125",
"20100126",
"20100127",
"20100128",
"20100129",
"20100201",
"20100202",
"20100203",
"20100204",
"20100205",
"20100208",
"20100209",
"20100210",
"20100211",
"20100212",
"20100215",
"20100216",
"20100217",
"20100218",
"20100219",
"20100224",
"20100225",
"20100226",
"20100227",
"20100301",
"20100302",
"20100303",
"20100304",
"20100305",
"20100309",
"20100310",
"20100311",
"20100312",
"20100315",
"20100316",
"20100317",
"20100318",
"20100319",
"20100322",
"20100323",
"20100324",
"20100325",
"20100326",
"20100329",
"20100330",
"20100331",
"20100401",
"20100402",
"20100405",
"20100406",
"20100407",
"20100408",
"20100409",
"20100412",
"20100413",
"20100414",
"20100415",
"20100416",
"20100419",
"20100420",
"20100421",
"20100422",
"20100423",
"20100426",
"20100427",
"20100428",
"20100429",
"20100430",
"20100504",
"20100505",
"20100506",
"20100507",
"20100511",
"20100512",
"20100513",
"20100514",
"20100517",
"20100518",
"20100519",
"20100520",
"20100521",
"20100524",
"20100525",
"20100526",
"20100527",
"20100528",
"20100531",
"20100601",
"20100602",
"20100603",
"20100604",
"20100607",
"20100608",
"20100609",
"20100610",
"20100611",
"20100615",
"20100616",
"20100617",
"20100618",
"20100621",
"20100622",
"20100623",
"20100624",
"20100625",
"20100628",
"20100629",
"20100630",
"20100701",
"20100702",
"20100705",
"20100706",
"20100707",
"20100708",
"20100709",
"20100712",
"20100713",
"20100714",
"20100715",
"20100716",
"20100719",
"20100720",
"20100721",
"20100722",
"20100723",
"20100726",
"20100727",
"20100728",
"20100729",
"20100730",
"20100802",
"20100803",
"20100804",
"20100805",
"20100806",
"20100809",
"20100810",
"20100811",
"20100812",
"20100813",
"20100816",
"20100817",
"20100818",
"20100819",
"20100820",
"20100823",
"20100824",
"20100825",
"20100826",
"20100827",
"20100830",
"20100831",
"20100901",
"20100902",
"20100903",
"20100906",
"20100907",
"20100908",
"20100909",
"20100910",
"20100913",
"20100914",
"20100915",
"20100916",
"20100917",
"20100920",
"20100921",
"20100922",
"20100923",
"20100924",
"20100927",
"20100928",
"20100929",
"20100930",
"20101001",
"20101004",
"20101005",
"20101006",
"20101007",
"20101008",
"20101011",
"20101012",
"20101013",
"20101014",
"20101015",
"20101018",
"20101019",
"20101020",
"20101021",
"20101022",
"20101025",
"20101026",
"20101027",
"20101028",
"20101029",
"20101101",
"20101102",
"20101103",
"20101108",
"20101109",
"20101110",
"20101111",
"20101112",
"20101113",
"20101115",
"20101116",
"20101117",
"20101118",
"20101119",
"20101122",
"20101123",
"20101124",
"20101125",
"20101126",
"20101129",
"20101130",
"20101201",
"20101202",
"20101203",
"20101206",
"20101207",
"20101208",
"20101209",
"20101210",
"20101213",
"20101214",
"20101215",
"20101216",
"20101217",
"20101220",
"20101221",
"20101222",
"20101223",
"20101224",
"20101227",
"20101228",
"20101229",
"20101230",
	}

	cmp := []string{
"20101008",
"20101011",
"20101012",
"20101013",
"20101014",
"20101015",
"20101018",
"20101019",
"20101020",
"20101021",
"20101022",
"20101025",
"20101026",
"20101027",
"20101028",
"20101029",
"20101101",
"20101102",
"20101103",
"20101108",
"20101109",
"20101110",
"20101111",
"20101112",
"20101113",
"20101115",
"20101116",
"20101117",
"20101118",
"20101119",
"20101122",
"20101123",
"20101124",
"20101125",
"20101126",
"20101129",
"20101130",
"20101201",
"20101202",
"20101203",
"20101206",
"20101207",
"20101208",
"20101209",
"20101210",
"20101213",
"20101214",
"20101215",
"20101216",
"20101217",
"20101220",
"20101221",
"20101222",
"20101223",
"20101224",
"20101227",
"20101228",
"20101229",
"20101230",
	}


	j := 0 // указатель сравниваемого массива дат
	for _,m := range	model {
		modelDate,_	:= toDate(m)
		//cmpDate,_ := toDate("20101230")
		//fmt.Println(i," compare:", modelDate.Equal(cmpDate))
		//func (t Time) Before(u Time) bool
		//Before reports whether the time instant t is before u. 
		curDate,_ := toDate(cmp[j])
		fmt.Printf("%v ",modelDate)
		if modelDate.Before(curDate) {
			fmt.Println("price: 0.0")
		} else {
			fmt.Println("price: VALUE")
			// получить слудующую дату
			j++
		}
	}

}

