package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	d   string
	txt string
)

func init() {
	flag.StringVar(&d, "d", "./down-data", "下载的文件夹目录，默认为当前文件夹下的down-data目录\n")
	flag.StringVar(&txt, "txt", "./url.txt", "下载的文件txt列表，默认为当前文件夹下的url.txt\n")
}

// 逐行读取文件内容
func ReadLines(fpath string) []string {
	fd, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var lines []string
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

func down(dir string, s string) string {
	//可以过滤url使其符合标准的url路径
	//	s = s[1:]
	//	s = s[:len(s)-1]
	u, err := url.Parse(s)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	filename := u.Path[1:]

	fpath := fmt.Sprintf(dir+"/%s", filename)
	newFile, err := os.Create(fpath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer newFile.Close()

	fmt.Println(filename + ":文件下载中...")
	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(s)
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return filename
}

func main() {
	flag.Parse()
	dir := d
	txt := txt

	urlList := ReadLines(txt)

	ch := make(chan string)
	for _, u := range urlList {
		go func(u string) {
			ch <- down(dir, u)
		}(u)
	}

	for i := 0; i < len(urlList); i++ {
		select {
		case result := <-ch:
			fmt.Println(result + "文件下载完成")
		case <-time.After(900 * time.Second):
			fmt.Println("Timeout..")
		}
	}
}
