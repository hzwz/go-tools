package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func worker(id int, jobs <-chan string) {
	for j := range jobs {
		getTitle(j)
	}
}

func getTitle(url string) {
	resp, err := http.Get(url + "/.git/head")
	if err != nil {
		errors.New("http请求失败")
		return
	}
	defer resp.Body.Close()
	/*
		if resp.StatusCode == 200 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			if find := strings.Contains(string(data), "refs/head"); find {
				fmt.Println("worker:", url+"/.git/")
			}

		}
	*/
	if resp.StatusCode == 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		reg := regexp.MustCompile(`<title>(.*?)</title>`)
		title := reg.FindStringSubmatch(string(data))
		if len(title) > 1 {
			fmt.Println(title[1])
		}
	}

}

func main() {
	jobs := make(chan string, 100)
	num := 30

	// 开启goroutine进行任务处理
	for w := 1; w <= num; w++ {
		go worker(w, jobs)
	}

	// 添加5个任务到任务队列
	//for j := 1; j <= 500; j++ {
	//jobs <- j
	//}
	//close(jobs)

	file, err := os.Open("gov.tw.txt")
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(file)

	lines := strings.Split(string(content), "\n")
	for _, url := range lines {
		url = strings.TrimSpace(url)
		jobs <- url

	}
	close(jobs)

}
