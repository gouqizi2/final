  package main
  
  import (
  	"net/http"
        "io/ioutil"
	"html/template"
  )

  type page struct {
  	title	string
  	body	[]byte
  }

  func (p *page) save() {
  	filename := p.title + ".txt"
  	ioutil.WriteFile(filename, p.body, 0600)
  }

  const lenPath = len("/view/")
   func loadPage(title string) *page {
  	filename := title + ".txt"
  	body ,_ := ioutil.ReadFile(filename)
  	return &page{title: title, body: body}
  } 
  func viewHandler(w http.ResponseWriter, r *http.Request) {
  	title := r.URL.Path[lenPath:]
  	p := loadPage(title)
  	renderTemplate(w, "view", p)
  }
   func saveHandler(w http.ResponseWriter, r *http.Request) {
  	title := r.URL.Path[lenPath:]
  	body := r.FormValue("username")+r.FormValue("password")
  	p := &page{title: title, body: []byte(body)}
  	p.save()
  	http.Redirect(w, r, "/view/"+title, http.StatusFound)
  }
   func login(w http.ResponseWriter, r *http.Request) {

   //     title := r.URL.Path[lenPath:]
	
        body := r.FormValue("username")+r.FormValue("password")
	b, _ := (ioutil.ReadFile("account.txt"))
	b2 := string(b[:])
	if b2 == body {
		http.Redirect(w, r, "/view/success/", http.StatusFound)
	}else {
		http.Redirect(w, r, "/view/fail/", http.StatusFound)
	}
  }
  func loginHandler(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        p := loadPage(title)
        renderTemplate(w, "login", p)
  } 
  func createHandler(w http.ResponseWriter, r *http.Request) {
  	title := r.URL.Path[lenPath:]
  	p := loadPage(title)
  	renderTemplate(w, "create", p)
  }
  func succHandler(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        p := loadPage(title)
        renderTemplate(w, "successful", p)
  }
  func failHandler(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        p := loadPage(title)
        renderTemplate(w, "failed", p)
  }  
  func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
  	t, _ := template.ParseFiles(tmpl+".html")
  	t.Execute(w,p)
  }


  func main() {
  	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/view/success/", succHandler)
	http.HandleFunc("/view/fail/", failHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/save/", saveHandler)
        http.HandleFunc("/login/", loginHandler)
        http.HandleFunc("/login/lg", login)
  	http.ListenAndServe(":8080", nil)
  }
