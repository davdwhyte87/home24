package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./*.html",
	}

	rnd = renderer.New(opts)
}
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/get_url_data", getUrlData).Methods("POST")
	port := ":300"
	http.ListenAndServe(port, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "home", GetURLDataReq{Url: "hooogo "})
}

type GetURLDataReq struct {
	Url string
}


type GetUrlResp struct {
	Url string
	HasLoginForm bool
	WorkingLinks int
	NonWorkingLinks int
	InternalLinks int
	ExternalLinks int
}
func getUrlData(w http.ResponseWriter, r *http.Request) {
	// var req GetURLDataReq
	r.ParseForm()
	url := r.FormValue("url")
	if url == ""{
		rnd.HTML(w, http.StatusOK, "home", GetURLDataReq{Url: ""})
		return
	}
	result:= dataProcess(url)
	result.Url = url
	
	rnd.HTML(w, http.StatusOK, "result", result)
}
func dataProcess(inputUrl string) GetUrlResp {
	// https://www.cantorsparadise.com/learn-anything-faster-by-using-the-feynman-technique-6565a9f7eda7
	// https://betterproposals.io/2/login/
	
	urlPersed, _ := url.Parse(inputUrl)
	pageHostName := urlPersed.Hostname()
	response, error := http.Get(inputUrl)
	if error != nil {
		println(error)
	}
	doc, error := goquery.NewDocumentFromReader(response.Body)
	if error != nil {
		println(error)
	}

	// check if a form exists on the page
	hasLoginForm := HasLoginOrSignup(doc)

	// h2s := doc.Find("h2")
	// title := doc.Find("title")
	workingLinks := 0
	nonWorkingLinks := 0
	internalLinks := 0
	externalLink := 0
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		// print (item.Is("a"))
		// println(strings.TrimSpace(item.Text()))

		link, ok := item.Attr("href")
		linkPersed, _ := url.Parse(link)
		linkHost := linkPersed.Hostname()
		// println("host name0", linkHost)
		if linkHost == pageHostName {
			internalLinks++
		} else {
			externalLink++
		}

		if ok == true {
			resp, hitErr := http.Get(link)
			if hitErr == nil {
				// println(link, resp.StatusCode)
				if resp.StatusCode > 400 {
					nonWorkingLinks++
				} else {
					workingLinks++
				}
			}
		}

		// println(link, resp.StatusCode)

	})

	println("------------------------")
	println("working links ", workingLinks)
	println("non working links", nonWorkingLinks)
	println("Internal link", internalLinks)
	println("External link", externalLink)

	return GetUrlResp{
		Url: "",
		WorkingLinks: workingLinks,
		NonWorkingLinks: nonWorkingLinks,
		InternalLinks: internalLinks,
		ExternalLinks: externalLink,
		HasLoginForm:hasLoginForm ,
	}
}

func HasLoginOrSignup(doc *goquery.Document) bool{
	// noOfForms := doc.Find("form").Size()
	isLoginPage := false
	doc.Find("form").Children().Each(func(index int, item *goquery.Selection) {
		// println(item.Html())
		// println(item.Attr("type"))
		itemAttr, _ := item.Attr("type")
		value, _ := item.Attr("value")
		if itemAttr == "submit" {
			// println("Button found")
			if strings.ToLower(value) == "sign in" || strings.ToLower(value) == "login" {
				isLoginPage = true
			}
		}
	})

	// attr, _:=doc.Find("form").Children().Find("input").Attr("type")
	// println(attr)

	println("Form data", isLoginPage)
	return isLoginPage
}
