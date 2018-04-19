  package main
  
  import (
  	"net/http"
        "io/ioutil"
	"html/template"
	"log"
	"net/rpc"
	"net"
  )
  type Req int

  func (s *Req) Process(title string, w *http.ResponseWriter) error {
	p := loadPage(title)
        renderTemplate(*w,"view",p)
	return nil
  }

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
  	rpc.Register(new(Req))
	rpc.HandleHTTP()
	l, e :=  net.Listen("tcp", ":1234")
	if e != nil {
        	log.Fatal("listen error:", e)
    	}
    	http.Serve(l, nil)
  }
