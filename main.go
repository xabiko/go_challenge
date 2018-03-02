package main

import (
  "html/template"
  "log"
  "fmt"
  "net/http"
  "os"
)

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
	t.Execute(w, r.URL.Path[1:])
}

func upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
        file, _, _ := r.FormFile("file_to_upload")
	fmt.Println(file)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload/", upload)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
