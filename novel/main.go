package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	getURL   = flag.Bool("c", false, "get novel url to a file.")
	getNovel = flag.Bool("n", false, "get novel content from files.")
	modNovel = flag.Bool("m", false, "modify novel title.")
)

func main() {
	chapterURL := config.URL + config.Chapter

	flag.Parse()
	flag.Usage = myUsage

	if *getURL == false && *getNovel == false && *modNovel == false {
		fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *getURL {
		println("start to download")
		writeNovelToFile(chapterURL, config.URLFile)
	}

	if *getNovel {
		downloadNovel(config.URLFile, config.NovelFile, config.StartPage, GetNovelByParam)
	}

	if *modNovel {
		handleText(config.NovelFile, "tmp.txt")
	}
}

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func writeNovelToFile(url, filename string) {
	pages := GetPages(url)
	var str string
	for _, page := range pages {
		str += fmt.Sprintf("[%s %s/%s]\n", page.Title, url, page.URL)
	}
	WriteToFile(filename, str)
}

func downloadNovel(src, dest string, firstChapter int, f func(string) []Novel) {
	outFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Cannot open file %s", dest)
	}

	bWriter := bufio.NewWriter(outFile)
	pages := getNovelFromFile(src, firstChapter)
	for _, page := range pages {
		novel := f(page.URL)[0]
		novel.Title = page.Title
		novelText := novel.Title + "\n\n    " + strings.TrimSpace(novel.Content) + "\n\n"
		bWriter.WriteString(novelText)
	}
	bWriter.Flush()
	outFile.Close()
}

func getNovelFromFile(filename string, firstChapter int) (pages []Page) {
	novelStr := ReadFromFile(filename)
	for i, novelLine := range strings.Split(novelStr, "\n") {
		if novelLine == "" {
			break
		}
		novelLine = strings.TrimPrefix(novelLine, "[")
		novelLine = strings.TrimRight(novelLine, "]")
		log.Println(novelLine)
		novelSplit := strings.Split(novelLine, " ")
		chapter := fmt.Sprintf("第%d章 %s", i+firstChapter, novelSplit[1])
		page := Page{
			Title: chapter,
			URL:   novelSplit[2],
		}
		pages = append(pages, page)
	}
	return pages
}

func handleText(inputFileName, outputFileName string) error {
	var re *regexp.Regexp
	var count int = 0

	inFile, err := os.Open(inputFileName)
	if err != nil {
		log.Printf("Cannot open text file: %s, err: [%v]", inputFileName, err)
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Cannot open file %s", outputFileName)
	}

	bWriter := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		// 更改章节
		re = regexp.MustCompile("^[\u7b2c](.*)[\u7ae0] ")
		matched := re.MatchString(line)
		if matched {
			println(line)
			replaceString := fmt.Sprintf("第%d章 ", config.StartPage+count)
			line = re.ReplaceAllString(line, replaceString)
			count++
		}
		bWriter.WriteString(line + "\n")
	}
	bWriter.Flush()

	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: %s, err: [%v]", inputFileName, err)
		return err
	}
	return nil
}
