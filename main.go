package main

import (
  "net/http"
  "encoding/json"
  "html/template"
  "strconv"
  "log"
  "os"
)

// JSON sample data to Go struct
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

// Only handler
func home(w http.ResponseWriter, r *http.Request) {

  t, _ := template.ParseFiles("index.html") //Link index.html

  switch r.Method {
  case "GET":
	t.Execute(w, template.HTML("Upload your JSON!"))  //Landing page will display initial message
  case "POST":
	// In-memory struct
	dat := Data{}
	// Read file from request
	file, _, err := r.FormFile("file_to_upload")
	if file == nil {
	  t.Execute(w, template.HTML("<pre>Error: No file chosen</pre>"))
	  return
	}
	defer file.Close()

	// Decode JSON into struct 
	err = json.NewDecoder(file).Decode(&dat)
	if err != nil {
	  t.Execute(w, template.HTML("<pre>Error: "+err.Error()+"</pre>"))
	  return
	}

	// Forging HTML response as a long string
	var total_superstring string = "<thead><tr><th>Job</th><th>Applicant Name</th><th>Email Address</th><th>Website</th><th>Skills</th><th>Cover Letter Paragraph</th></tr></thead><tbody>"

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

			desc := "<td rowspan="+strconv.Itoa(app_rowspan)+" class=applicant-name>"+dat.Applicants[app].Name+"</td><td rowspan="+strconv.Itoa(app_rowspan)+"><a href=mailto:"+dat.Applicants[app].Email+">"+dat.Applicants[app].Email+"</a></td><td rowspan="+strconv.Itoa(app_rowspan)+"><a href=http://"+dat.Applicants[app].Website+"/>"+dat.Applicants[app].Website+"</a></td>"+s[0]+"<td rowspan="+strconv.Itoa(app_rowspan)+">"+dat.Applicants[app].CoverLetter+"</td></tr>"

			for i:=1; i<len(s); i++ {desc += "<tr>"+s[i]+"</tr>"}
			a = append(a, desc)
			job_rowspan += app_rowspan
		}
		total_superstring += "<tr><td rowspan="+strconv.Itoa(job_rowspan)+" class=&#34;job-name&#34;>"+dat.Jobs[job].Name+"</td>"+a[0]
		for i:=1; i<len(a); i++ {total_superstring += "<tr>"+a[i]}
	}

	unique_skills := []string {dat.Skills[0].Name}  //Custom subsetting of unique values
	for i := range dat.Skills {
		l := len(unique_skills)
		for j:=0; j<l; j++{
			if dat.Skills[i].Name == unique_skills[j] {break}
			if j==(l-1) {
			  unique_skills = append(unique_skills, dat.Skills[i].Name)
			}
		}
	}

	// End of HTML response
	total_superstring += "</tbody><tfoot><tr><td colspan=6>"+strconv.Itoa(len(dat.Applicants))+" Applicants, "+strconv.Itoa(len(unique_skills))+" Unique Skills</tr></tfoot>"

	t.Execute(w, template.HTML(total_superstring))
  }
}

func main() {
	http.HandleFunc("/", home) //Only endpoint
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
