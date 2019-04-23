//
// загрузчик данных с сайта Финам
// по мотивам статьи https://habr.com/ru/post/332700/
//
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
					// цифры
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

func transform(rootDir string) {
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
			fmt.Printf("visited file: %q\n", path)
			// выбрать из файла котировок(quotations) данные тикер, дата, цена закрытия 
			// и поместить их в буфер для последующего формирвоания сводной таблицы
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Unable to open file:", path, err)
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				// получить поля
				// tiker:field[0],date:field[2],close:field[7]
				fields := strings.Split(scanner.Text(), ";")
				//fmt.Println("*",fields)
				// преобразовать в float64
				floatNum,_ := toFloat64(fields[7])
				fmt.Printf("tiker: %s date: %q  price: %q[%.4f]\n",fields[0],fields[2],fields[7],floatNum)
			}
		} // eof if-else

		return nil
	})

	if err != nil {
		//fmt.Printf("error walking the path %q: %v\n", tmpDir, err)
		fmt.Println("transform error:",err)
		return
	}
}


// convert finam quotations' prices value to Golang float64
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

	//toFloatEx()

	//getAssetParams()
	//downloadAssetHistory()

	//downloadEx()
	//dirEx()
	//mapEx()

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

