package common

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetFileList read directory path and return file list
func GetFileList(dir string) []string {
	files := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if strings.Contains(path, ".git") {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		baseName := filepath.Base(path)
		if !strings.Contains(baseName, "-") {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

// ParseNews get news from file and return news list
func ParseNews(path string) (newsList []News, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`(?i)#+\s*?GoCN\D*([\d-]*).*\n+\D*((.*\n)+?\n)`)
	allMatch := reg.FindAllSubmatch(data, -1)
	for _, item := range allMatch {
		loc, _ := time.LoadLocation("Local")
		ctime, err := time.ParseInLocation("2006-01-02", string(item[1]), loc)
		if err != nil {
			continue
		}
		subReg := regexp.MustCompile(`(\d)\.\s*(.*)\s*(http.*)\n?`)
		subAll := subReg.FindAllSubmatch(item[2], -1)
		for _, subItem := range subAll {
			id, err := strconv.ParseInt(string(subItem[1]), 10, 64)
			if err != nil {
				continue
			}
			title := strings.TrimSpace(string(subItem[2]))
			sURL, err := url.Parse(strings.TrimSpace(string(subItem[3])))
			if err != nil {
				continue
			}
			singleNews := News{id, title, *sURL, ctime}
			newsList = append(newsList, singleNews)
		}
	}
	return
}
