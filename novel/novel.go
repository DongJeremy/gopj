package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Page info
type Page struct {
	Title string
	URL   string
}

// Novel info
type Novel struct {
	Title   string
	Content string
}

// NovelURL send php request
type NovelURL struct {
	id  string
	qid string
}

// GetPages 获取分页
func GetPages(url string) []Page {
	doc, err := Curl(url)
	if err != nil {
		log.Fatal(err)
	}

	return ParsePages(doc)
}

// GetNovel 获取novel
func GetNovel(url string) []Novel {
	doc, err := Curl(url)
	if err != nil {
		log.Fatal(err)
	}

	return ParseNovel(doc)
}

// Split split sep by / or .
func Split(r rune) bool {
	return r == '/' || r == '.'
}

// GetNovelByParam 通过 param 获取novel
func GetNovelByParam(urlString string) (novels []Novel) {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	tmpArr := strings.FieldsFunc(u.Path, Split)
	BaseFMT := config.URL + config.PhpURL
	baseURL := fmt.Sprintf(BaseFMT, tmpArr[2], tmpArr[3])
	println(baseURL)

	resText, err := CurlText(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	//resText = DecodeGBK(resText)
	resText = strings.Replace(resText, `&nbsp;&nbsp;&nbsp;&nbsp;`, `    `, -1)
	var re *regexp.Regexp
	//去除br
	re, _ = regexp.Compile("<br>")
	resText = re.ReplaceAllString(resText, "")
	re, _ = regexp.Compile("<br />")
	resText = re.ReplaceAllString(resText, "\n")

	//去除viewall
	re, _ = regexp.Compile("<font[\\S\\s]+?</font>")
	resText = re.ReplaceAllString(resText, "")

	novel := Novel{
		Title:   "",
		Content: resText,
	}

	novels = append(novels, novel)
	return
}

// ParsePages 分析分页
func ParsePages(doc *goquery.Document) (pages []Page) {
	doc.Find("#list > dl > dd > a").Each(func(i int, s *goquery.Selection) {
		page := DecodeGBK(s.Text())
		url, _ := s.Attr("href")

		pages = append(pages, Page{
			Title: page,
			URL:   url,
		})
	})

	return pages
}

// ParseNovel 分析novel数据
func ParseNovel(doc *goquery.Document) (novels []Novel) {
	doc.Find("#BookCon").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h1").Text()
		title = DecodeGBK(title)

		subtitle := strings.Split(title, " ")
		content := DecodeGBK(s.Find("#pagecontent").Text())
		content = strings.Replace(content, `聽聽聽聽`, `    `, -1)
		novel := Novel{
			Title:   subtitle[0] + " " + subtitle[1],
			Content: content,
		}

		novels = append(novels, novel)
	})

	return novels
}
