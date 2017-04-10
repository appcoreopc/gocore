package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var image []byte

type CoreFile struct {
}

func (corefile CoreFile) GetFileContent(filename string) []byte {

	var err error
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return content
}

func handleFile(w http.ResponseWriter, r *http.Request) {

	nbytes, nchunks := int64(0), int64(0)

	_ = nbytes
	_ = nchunks

	websocket, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	conn, bwriter, okhijack := websocket.Hijack()

	if okhijack != nil {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
	}

	defer conn.Close()

	/*
		bwriter.WriteString("lets go man!!!")

		time.Sleep(1 * time.Second)
		time.Sleep(1 * time.Second)
		bwriter.WriteString("lets go man!!!" + time.Now().String())
		time.Sleep(1 * time.Second)
		time.Sleep(1 * time.Second)
		bwriter.WriteString("lets go man!!!" + time.Now().String())
		time.Sleep(1 * time.Second)

		println("done and flushing out contents!!")

		bwriter.Flush()
	*/

	/*
		for i := 0; i < 4; i++ {
			time.Sleep(1 * time.Second)
			bwriter.WriteString("lets go man!!!" + time.Now().String())
			bwriter.Flush() // you need this, otherwise no output
		}

	*/

	file, fopenOk := os.Open("app.txt")

	if fopenOk != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
	}

	println("reading file")

	reader := bufio.NewReader(file)
	buff := make([]byte, 1024)

	for {

		n, e := reader.Read(buff[:cap(buff)])
		buff := buff[:n]

		if n != 0 {
			if e == nil {
				c := string(buff)
				bwriter.WriteString(c)
				bwriter.Flush()
				println("flusing" + c)
				time.Sleep(5 * time.Second)
				continue
			}
			if e == io.EOF {
				break
			}
		} else if n == 0 {
			println("statis mode!")
			time.Sleep(5 * time.Second)
		}

		nchunks++
		nbytes += int64(len(buff))

		if e != nil && e != io.EOF {
			log.Fatal(e)
		}
	}

}

func main() {
	http.HandleFunc("/", handleFile)
	fmt.Println("start http listening :8081")
	err := http.ListenAndServe(":8081", nil)
	fmt.Println(err)
}
