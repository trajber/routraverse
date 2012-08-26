package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileData chan []byte

func Walk(path string, c chan FileData) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range files {
		abs := fmt.Sprintf("%s%c%s", path, os.PathSeparator, v.Name())

		if v.IsDir() {
			Walk(abs, c)
			continue
		}

		x := make(chan []byte)
		go func(abs string) {
			data, _ := ioutil.ReadFile(abs)
			x <- data
			close(x)
		}(abs)

		c <- x
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <path>")
		os.Exit(-1)
	}

	ch := make(chan FileData)

	go func() {
		Walk(os.Args[1], ch)
		close(ch)
	}()

	for {
		c, ok := <-ch

		if ok == false {
			break
		}

		data := <-c
		fmt.Println(string(data))
	}
}
