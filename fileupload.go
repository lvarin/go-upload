package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var index string = " "

func fileUploadHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println("Error Getting File", err)
			return
		}
		defer file.Close()

		out, pathError := ioutil.TempFile("temp-images", "*-"+header.Filename)
		if pathError != nil {
			log.Println("Error Creating a file for writing", pathError)
			return
		}
		defer out.Close()

		_, copyError := io.Copy(out, file)
		if copyError != nil {
			log.Println("Error copying", copyError)
		}
		fmt.Fprintln(w, "File Uploaded Successfully! ")
		fmt.Fprintln(w, "Name of the File: ", header.Filename)
		fmt.Fprintln(w, "Size of the File: ", header.Size)
	}
}

func main() {
	log.Printf("Upload GO server\n")

	// Load the index.html file into `index` variable.
	b, err := ioutil.ReadFile("index.html") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	index = string(b)

	// k8strain.Test(os.Stdout)
	RestServer()
}

// RestServer registers all URLs and start the waiting.
func RestServer() {
	jobRouter := mux.NewRouter().StrictSlash(true)
	//replyRouter := mux.NewRouter().StrictSlash(true)

	jobRouter.HandleFunc("/", homePage())
	jobRouter.HandleFunc("/upload", fileUploadHandler())

	// Start main loops
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(8080),
			handlers.LoggingHandler(os.Stdout, jobRouter)))
		wg.Done()
	}()

	wg.Wait()

}

func homePage() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, index)
	}
}
