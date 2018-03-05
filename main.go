package main

import (
  "strconv"
  "encoding/json"
  //"io"
  //"bytes"
  "html/template"
  "log"
  "fmt"
  "net/http"
  "os"
)

//var data []byte

type Data struct {
	Applicants []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		Website     string `json:"website"`
		CoverLetter string `json:"cover_letter"`
		JobID       int    `json:"job_id"`
	} `json:"applicants"`
	Jobs []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"jobs"`
	Skills []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ApplicantID int    `json:"applicant_id"`
	} `json:"skills"`
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	__welcome := "<b>Upload your JSON</b>"
	t.Execute(w, template.HTML(__welcome))
}

func upload(w http.ResponseWriter, r *http.Request) {

	// Create memory buffer
	//var dat map[string][]interface{}
	dat := Data{}
	// Read file from request
	file, _, err := r.FormFile("file_to_upload")
	if err != nil {
        panic(err)
	}
	defer file.Close()
	// Copy the file data to my buffer
	err = json.NewDecoder(file).Decode(&dat)
	if err != nil {
	panic(err)
	}

	var total_superstring string

	for job := range dat.Jobs {
		var a []string
		job_rowspan := 0

		for app := range dat.Applicants {
			var s []string
			app_rowspan := 0
			if dat.Jobs[job].ID != dat.Applicants[app].JobID {continue}

			for ski := range dat.Skills {
				if dat.Applicants[app].ID != dat.Skills[ski].ApplicantID {continue}
				app_rowspan += 1
				s = append(s, "<td>"+dat.Skills[ski].Name)
			}

			desc := "<td rowspan="+strconv.Itoa(app_rowspan)+" class=applicant-name>"+dat.Applicants[app].Name+"</td><td rowspan="+strconv.Itoa(app_rowspan)+"><a href=&#34;mailto:"+dat.Applicants[app].Email+"&#34;>"+dat.Applicants[app].Email+"</a></td><td rowspan="+strconv.Itoa(app_rowspan)+"><a href=&#34;http://"+dat.Applicants[app].Website+"/&#34;>"+dat.Applicants[app].Website+"</a></td>"+s[0]+"<td rowspan="+strconv.Itoa(app_rowspan)+">"+dat.Applicants[app].CoverLetter+"&#34;</td></tr>"

			for i:=1; i<len(s); i++ {desc += "<tr>"+s[i]+"</tr>"}
			a = append(a, desc)
			job_rowspan += app_rowspan
		}
		total_superstring += "<tr><td rowspan="+strconv.Itoa(job_rowspan)+" class=&#34;job-name&#34;>"+dat.Jobs[job].Name+"</td>"+a[0]
		for i:=1; i<len(a); i++ {total_superstring += "<tr>"+a[i]}
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, template.HTML(total_superstring))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
