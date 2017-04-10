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
var defaultFileName string = "app.txt"

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

	file, fopenOk := os.Open(defaultFileName)

	if fopenOk != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
	}

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
				//time.Sleep(5 * time.Second)
				continue
			}
			if e == io.EOF {
				break
			}

		} else if n == 0 { // if no new content
			// putting on a wait time
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
