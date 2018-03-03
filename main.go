package main

import (
  "strings"
  "io"
  "bytes"
  "html/template"
  "log"
  "fmt"
  "net/http"
  "os"
)

var data []byte

type Jobs struct {
	People Applicant
}

type Applicant struct {
	Name string
	Email string
	Website string
	Skills []string
	CoverLetter string
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file_to_upload")
	if err != nil {
        panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct,
	// but this will work as an example
	contents := Buf.String()
	fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	return
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
