package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const NumLen = 4

var key int
var files []string
var param string
var strlen int

func getFiles(args []string) {
	for _, file := range args {
		if f, err := os.Stat(file); err == nil {
			if f.IsDir() {
				fmt.Println("ERROR: ", file, "is directory!")
				files = make([]string, 0)
				return
			}
		} else {
			fmt.Println("ERROR: ", file, "no such file or directory!")
			files = make([]string, 0)
			return
		}
		files = append(files, file)
	}
}

func inputCheck(args []string) bool {
	if len(args) < 2 {
		printError()
		return false
	}
	param = args[1]
	if param == "-d" || param == "-e" {
		if len(args[2]) > NumLen {
			printError()
			return false
		}
		num, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
			return false
		}
		if num <= 0 {
			printError()
			return false
		}
		key = num
		getFiles(args[3:])
		return true
	}
	if param == "-c" {
		getFiles(args[2:])
		return true
	}
	return true
}

func printError() {
	fmt.Println("\n  请重新输入,格式如下:")
	fmt.Println("\n  1输入需要插入到文件最前面的字节数:1-9999")
	fmt.Println("\n  2输入需要加密的文件(文件不需要和pt文件同目录): file1 file2 ... fileN")
	fmt.Println("\n  3例子: $ pt -e/-d 99 file1.jpg file2.jpg ... D:\\fileN.jpg")
}

func outputPaht(key int, file string) string {
	dirname := path.Dir(file)
	fullname := path.Base(file)
	filesuffix := path.Ext(fullname)
	filename := strings.TrimSuffix(fullname, filesuffix)
	return path.Join(dirname, filename+"_"+strconv.Itoa(key)+filesuffix)
}

func encryptImg(key int, file string) {
	if src, err := ioutil.ReadFile(file); err == nil {
		b := make([]byte, key)
		if _, err := rand.Read(b); err == nil {
			dest := append(b, src...)
			outfile := outputPaht(key, file)
			err := ioutil.WriteFile(outfile, dest, 0755)
			if err != nil {
				panic(err)
			}
		}
	}
}

func decryptImg(key int, file string) {
	if src, err := ioutil.ReadFile(file); err == nil {
		dest := src[key:]
		outfile := outputPaht(key, file)
		err := ioutil.WriteFile(outfile, dest, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func countStr(str string) {
	lenght := strings.Count(str, "") - 1
	strArr := strings.Split(str, ".")
	last := strings.Count(strArr[len(strArr)-1], "") - 1
	fmt.Println(
		"Length: ",
		strconv.Itoa(lenght-last),
		strconv.Itoa(last),
		strconv.Itoa(lenght),
		str)

}

func changeName(file string, n int) {
	filename := path.Base(file)
	filesuffix := path.Ext(filename)
	name := path.Join(path.Dir(file), "a"+strconv.Itoa(n)+filesuffix)
	err := os.Rename(file, name)
	if err != nil {
		panic(err)
	}
}

func main() {
	if inputCheck(os.Args) {
		switch param {
		case "-e":
			for _, file := range files {
				encryptImg(key, file)
			}
		case "-d":
			for _, file := range files {
				decryptImg(key, file)
			}
		case "-c":
			for n, file := range files {
				changeName(file, n)
			}
		default:
			countStr(param)
		}
	}
}
