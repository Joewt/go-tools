package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	d    string
	txt  string
	data string
)

type JsonStruct struct {
	Data []DirStruct `json:"data"`
}

type DirStruct struct {
	Dir      string `json:"dir"`
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

func init() {
	flag.StringVar(&d, "d", "./down-data", "下载的文件夹目录，默认为当前文件夹下的down-data目录\n")
	flag.StringVar(&txt, "txt", "./url.txt", "下载的文件txt列表，默认为当前文件夹下的url.txt\n")
	flag.StringVar(&data, "data", "./data.json", "下载的文件结构json，默认为当前文件夹下的data.json\n")
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
	client := http.Client{Timeout: 1800 * time.Second}
	resp, err := client.Get(s)
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return filename
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func downv2(data DirStruct) string {
	dir := data.Dir
	s := data.Url
	name := data.Filename
	if !isExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			//return "目录创建失败"
		}
	}
	fpath := fmt.Sprintf(dir + "/" + name)
	newFile, err := os.Create(fpath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer newFile.Close()

	fmt.Println(name + ":文件下载中...")
	client := http.Client{Timeout: 1800 * time.Second}
	resp, err := client.Get(s)
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "文件copy失败"
	}
	//wg.Done()
	return name
}

func readData(source string) ([]byte, error) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	flag.Parse()
	//v1  批量简单下载文件，文件内为一行行的url地址
	//	dir := d
	//	txt := txt
	/*
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
	*/

	//v2 根据json文件自动下载数据，目录关系json配置里定义好
	data := data

	d, err := readData(data)
	if err != nil {
		panic(err.Error())
	}
	dataList := JsonStruct{}
	err = json.Unmarshal(d, &dataList)
	if err != nil {
		panic(err.Error())
	}
	ch := make(chan string)
	//var wg sync.WaitGroup
	for _, v := range dataList.Data {
		//	wg.Add(1)
		go func(u DirStruct) {
			fmt.Println("处理:", u)
			ch <- downv2(u)
		}(v)
	}
	for {
		select {
		case result := <-ch:
			fmt.Println(result + "文件下载完成")
		case <-time.After(1800 * time.Second):
			fmt.Println("Tiemout...")
		}
	}
}
