package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var chnNumChar = map[string]int{
	"〇": 0,
	"零": 0,
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"七": 7,
	"八": 8,
	"九": 9,

	"壹": 1,
	"贰": 2,
	"叁": 3,
	"肆": 4,
	"伍": 5,
	"陆": 6,
	"柒": 7,
	"捌": 8,
	"玖": 9,

	"貮": 2,
	"两": 2,
}

var chnNumUnit = map[string]int{
	"十": 10,
	"拾": 10,
	"百": 100,
	"佰": 100,
	"千": 1000,
	"仟": 1000,
	"万": 10000,
	"萬": 10000,
	"亿": 100000000,
	"億": 100000000,
	"兆": 1000000000000,
}

// IsChineseChar 判断字符串是否包含中文字符
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// ChineseToNumber ch
func ChineseToNumber(chnStr string) int {
	var cnChar = strings.Split(chnStr, "")
	var unit = 0
	var cnDig = make([]string, 10)
	//如果不是汉字则返回
	if !IsChineseChar(chnStr) {
		xInt, _ := strconv.Atoi(chnStr)
		return xInt
	}
	for i := len(cnChar) - 1; i >= 0; i-- {
		if v, ok := chnNumUnit[cnChar[i]]; ok {
			unit = v
			if unit == 10000 {
				cnDig = append(cnDig, "w")
				unit = 1
			} else if unit == 100000000 {
				cnDig = append(cnDig, "y")
				unit = 1
			} else if unit == 1000000000000 {
				cnDig = append(cnDig, "z")
				unit = 1
			}
		} else {
			if v2, ok := chnNumChar[cnChar[i]]; ok {
				dig := v2
				if unit != 0 {
					dig = dig * unit
					unit = 0
				}
				cnDig = append(cnDig, strconv.Itoa(dig))
			}
		}
	}

	if unit == 10 {
		cnDig = append(cnDig, "10")
	}

	var tmp = 0
	var ret = 0

	for i := len(cnDig) - 1; i >= 0; i-- {
		if cnDig[i] == "w" {
			tmp *= 10000
			ret += tmp
			tmp = 0
		} else if cnDig[i] == "y" {
			tmp *= 100000000
			ret += tmp
			tmp = 0
		} else if cnDig[i] == "z" {
			tmp *= 1000000000000
			ret += tmp
			tmp = 0
		} else {
			xInt, _ := strconv.Atoi(cnDig[i])
			tmp += xInt
		}
	}

	ret += tmp
	return ret
}

// WriteToFile write info to files.
func WriteToFile(filename, content string) {
	if err := ioutil.WriteFile(filename, []byte(content), 0644); err != nil {
		fmt.Println(err)
	}
	fmt.Println("write file success...")
}

// ReadFromFile 读取到file中，再利用ioutil将file直接读取到[]byte中
func ReadFromFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("read file fail", err)
		return ""
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return ""
	}

	return string(fd)
}

// DecodeGBK decode
func DecodeGBK(src string) (dst string) {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewDecoder()))
	if err == nil {
		dst = string(data)
	}
	return
}

// Curl get url
func Curl(baseURL string) (doc *goquery.Document, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Referer", baseURL)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err = goquery.NewDocumentFromResponse(res)
	return
}

// CurlText get url
func CurlText(baseURL string) (resText string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Referer", baseURL)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	resText = string(body)
	return
}
