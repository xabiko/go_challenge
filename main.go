package main

import (
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

	r.ParseMultipartForm(32 << 20)
        file, _, err := r.FormFile("file_to_upload")
	if err != nil {
	   fmt.Println(err)
	   return
        }
        defer file.Close()
	file.Read(data)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload/", upload)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
