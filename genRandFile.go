package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

var (
	d string
	p int
	n int
	s float64
)

func init() {
	flag.StringVar(&d, "d", "./data", "文件夹目录,默认为当前文件下data目录\n")
	flag.IntVar(&p, "p", 1, "模式：\n1-往目录下所有文件追加随机字符,\n2-生成文件\n")
	flag.IntVar(&n, "n", 100, "生成的文件数，默认100\n")
	flag.Float64Var(&s, "s", 1, "生成的文件大小，默认1k\n")
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func RewriteFile(pathname string) {
	dir_list, e := ioutil.ReadDir(pathname)
	if e != nil {
		fmt.Println(e)
		return
	}
	if len(dir_list) == 0 {
		fmt.Println("文件为空!!")
		return
	}
	for _, file := range dir_list {
		if file.IsDir() {
			RewriteFile(pathname + "/" + file.Name())
		} else {
			s := UniqueId()

			filename := pathname + "/" + file.Name()
			fd, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			defer fd.Close()
			if _, err2 := io.WriteString(fd, s); err2 == nil {
				fmt.Println("success save file:" + pathname + "/" + file.Name())
			}
		}
	}
}

func createFile(size float64, filename string) {
	size = math.Ceil(size)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	count := math.Ceil(float64(size) / 1000)
	count_64 := int64(int(count))
	var i int64
	var length int
	for i = 0; i < count_64; i++ {
		if i == (count_64 - 1) {
			length = int(int64(size) - (i)*1000)
		} else {
			length = 1000
		}
		s := UniqueId()
		f.WriteAt([]byte(strings.Repeat(s, length)), i*1000)
	}
}

func GenFile(pathname string, filenum int, filesize float64) {
	_, err := os.Stat(pathname)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < filenum; i++ {
		s := UniqueId()

		filename := pathname + "/" + s
		createFile(filesize, filename)
		fmt.Println("success gen file filename:" + filename)
	}
}

func main() {

	flag.Parse()

	dir := d
	pattern := p
	filenum := n

	if pattern == 1 {
		RewriteFile(dir)
	} else if pattern == 2 {
		GenFile(dir, filenum, s*32*math.Ceil((s/32)))
	}
}
