package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"sync"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type arrayFlags []string

func main() {
	var headers arrayFlags
	filepath := flag.String("i", "urls.txt", "Path to file containing urls to test")
	flag.Var(&headers, "H" ,"Headers ex: -H='Cookie: PHPSESSID=shjhjdgvbhjhvnv'")
	flag.Parse()
	urls := getUrlsFromFile(string(*filepath))

	var wg sync.WaitGroup

	for i:=0;i<len(urls);i++{
		wg.Add(1)
		go func(i int, headers arrayFlags) {
			defer wg.Done()
			url := urls[i]
			fmt.Printf("[+] Checking %s\n", url)
			data, err := fetchURL(url ,headers)
			if err != nil{
				return
			}
			checkPostMessage(data, url)
		}(i, headers)
	}
	wg.Wait()
}

func fetchURL(url string, headers arrayFlags) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		return nil, err
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()

	req.WithContext(ctx)

	for _, header := range headers{
		header = strings.ReplaceAll(header, " ", "")
		h := strings.Split(header, ":")
		req.Header.Add(h[0], h[1])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes, nil
}

func checkPostMessage(bytes []byte, url string){
	body := string(bytes)
	lbody := strings.ToLower(body)
	abody := strings.Split(lbody, "\n")

	for i, line := range  abody{
		if strings.Contains(line, "addeventlistener(\"message"){
			fmt.Println(url)
			fmt.Printf("%d: postMessage event listener detected!\n", i+1)
			out := strings.Trim(line, " ")
			fmt.Printf(ErrorColor, out+"\n")
		}

		if strings.Contains(line, "addeventlistener('message"){
			fmt.Println(url)
			fmt.Printf("%d: postMessage event listener detected!\n", i+1)
			out := strings.Trim(line, " ")
			fmt.Printf(ErrorColor, out+"\n")
		}

		if strings.Contains(line, "window.attachevent(\"message"){
			fmt.Println(url)
			fmt.Printf("%d: postMessage event listener detected!\n", i+1)
			out := strings.Trim(line, " ")
			fmt.Printf(ErrorColor, out+"\n")
		}

		if strings.Contains(line, "window.attachevent('message"){
			fmt.Println(url)
			fmt.Printf("%d: postMessage event listener detected!\n", i+1)
			out := strings.Trim(line, " ")
			fmt.Printf(ErrorColor, out+"\n")
		}

		if strings.Contains(line, "onmessage"){
			fmt.Println(url)
			fmt.Printf("%d: postMessage function detected!\n", i+1)
			out := strings.Trim(line, " ")
			fmt.Printf(ErrorColor, out+"\n")
		}
	}
}

func getUrlsFromFile(path string) []string {
	var urls []string
	file, err := os.Open(path)
	if err != nil{
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		urls = append(urls, scanner.Text())
	}
	return urls
}

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}