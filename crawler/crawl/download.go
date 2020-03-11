package crawl

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getChapterListURL(chapterIndex string) {
	doc, err := goquery.NewDocument(chapterIndex)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	doc.Find("#list dl dd").Each(func(i int, s *goquery.Selection) {
		text, _ := s.Attr("value")
		listURL := text
		ch1 <- listURL
	})
}

func getChapterNumber(suffix string) {
	for listURL := range ch1 {
		doc, err := goquery.NewDocument(listURL)
		if err != nil {
			fmt.Println("err:", err)
			os.Exit(1)
		}

		doc.Find(".idx-ol li").Each(func(i int, s *goquery.Selection) {
			text, _ := s.Find("a").Attr("href")
			//去掉最左边的'/'
			u, err := url.Parse(text)
			if err != nil {
				log.Println("Parse url failed:", text, err)
				return
			}
			//tmp := strings.TrimLeft(u.Path, "/292926/")
			tmpString := u.Path
			tmpString = string(tmpString[8:])
			println(tmpString)
			tmpString = strings.TrimRight(tmpString, suffix)
			//listUrl := "https://m.zhulang.com/BdRead/index/book_id/292926/ch_id/" + tmpString + "/channel/200001.html"
			ch2 <- listURL
		})
	}
}

func downloadNovel() {
	for novelListURL := range ch2 {
		saveNovel(novelListURL)
	}
}

func saveNovel(nURL string) {
	log.Println(nURL)
	suffix := ".html"
	u, err := url.Parse(nURL)
	if err != nil {
		log.Println("Parse url failed:", nURL, err)
		return
	}
	//去掉最左边的'/'
	tmp := strings.TrimLeft(u.Path, "/html/")
	prefix := strings.TrimRight(tmp, suffix)
	filename := novelDir + string(os.PathSeparator) + strings.ToLower(strings.Replace(prefix, "/", "-", -1)) + ".txt"

	exists := isExist(filename)
	if exists {
		return
	}

	nFile, err := os.Create(filename)
	if err != nil {
		log.Println("create file failed:", filename, err)
		return
	}

	content := parseNovelContent(nURL)
	ch3 <- 1
	defer nFile.Close()
	nFile.WriteString(content)
}

//解析novel url
func parseNovelContent(listURL string) string {
	doc, _ := goquery.NewDocument(listURL)
	println("url:", listURL)
	title := doc.Find(".wrap .p-tit h2").Eq(0).Text()
	//println(title)
	textStr, _ := doc.Find(".wrap .art").Html()
	//println(textStr)
	return convertContent(title, textStr)
}

func convertContent(title, content string) string {
	tmpStr := strings.Replace(content, `聽聽聽聽`, `    `, -1)
	var re *regexp.Regexp

	//去除br
	re, _ = regexp.Compile("<p>")
	tmpStr = re.ReplaceAllString(tmpStr, "\n    ")
	re, _ = regexp.Compile("</p>")
	tmpStr = re.ReplaceAllString(tmpStr, "\n")
	re, _ = regexp.Compile("<p class=\"tc\">......")
	tmpStr = re.ReplaceAllString(tmpStr, "")

	//去除viewall
	re, _ = regexp.Compile("<!--div[\\S\\s]+?</div-->")
	tmpStr = re.ReplaceAllString(tmpStr, "")

	//章节替换
	// re, _ = regexp.Compile("第([\\S\\s]+?)章([\\S\\s]*)")
	// if re.MatchString(title) {
	// 	result := re.FindStringSubmatch(title)
	// 	title = fmt.Sprintf("第%d章 %s", ChineseToNumber(result[1]), result[2])
	// }

	return title + "\n\n    " + strings.TrimSpace(tmpStr) + "\n\n"
}
