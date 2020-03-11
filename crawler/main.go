package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//crawl.DownloadNovels()
	GetChapterListURL("http://www.biquger.com/biquge/39649/")
}

// GetChapterListURL d
func GetChapterListURL(chapterIndex string) {
	doc, err := goquery.NewDocument(chapterIndex)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	fmt.Println(doc)
	doc.Find("#list dl dd").Each(func(i int, s *goquery.Selection) {
		text, _ := s.Find("a").Attr("href")
		listURL := text
		fmt.Println(listURL)
	})
}
