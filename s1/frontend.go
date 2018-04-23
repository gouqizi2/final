  package main
  
  import (
  	"net/http"
        "io/ioutil"
	"html/template"
	"log"
	"net/rpc"
//	"net"
  )
var curUser = "Guest"
  type Page struct {
	Title string
	Body    []byte
  }
  type Pass struct {
	Username string
	Page    Page
  }
  func (p *Page) save() {
  	filename := p.Title + ".txt"
  	ioutil.WriteFile(filename, p.Body, 0600)
  }

  type Info struct {
        Username string
        Password string
  }
  const lenPath = len("/view/")
   func loadPage(title string) *Page {
  	filename := title + ".txt"
  	body ,_ := ioutil.ReadFile(filename)
  	return &Page{Title: title, Body: body}
  } 
  func viewHandler(w http.ResponseWriter, r *http.Request) {
  	title := r.URL.Path[lenPath:]
  	p := loadPage(title)
  	renderTemplate(w, "view", p)
  }
   func saveHandler(w http.ResponseWriter, r *http.Request) {
  	title := r.URL.Path[lenPath:]
//  	body := r.FormValue("username")+r.FormValue("password")
  	i := &Info{Username: r.FormValue("username"), Password: r.FormValue("password")}
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
        	log.Fatal("dialing:", err)
	}

  	var args = i
	var reply bool
	err = client.Call("Req.Create", args, &reply)
	if err != nil {
        	log.Fatal("arith error:", err)
	}
  	http.Redirect(w, r, "/view/"+title, http.StatusFound)
  }
   func login(w http.ResponseWriter, r *http.Request) {	
//        body := r.FormValue("username")+r.FormValue("password")
	U := r.FormValue("username")
	P := r.FormValue("password")
	i := &Info{Username: U, Password: P}
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
        	log.Fatal("dialing:", err)
	}

	var args = i
	var reply bool
	err = client.Call("Req.Login",args, &reply)
	if err != nil {
        	log.Fatal("arith error:", err)
	}
	if reply == true {
		curUser = U
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
  func postHandler(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        p := loadPage(title)
        renderTemplate(w, "post", p)
  }
  func po(w http.ResponseWriter, r *http.Request) {
//        body := r.FormValue("username")+r.FormValue("password")
        p := &Page{Title: r.FormValue("title"), Body: []byte(r.FormValue("content"))}
	p2 := &Pass{Username: curUser, Page: *p}
        client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
        if err != nil {
                log.Fatal("dialing:", err)
        }

        var args = p2
        var reply bool
        err = client.Call("Req.Post",args, &reply)
        if err != nil {
                log.Fatal("arith error:", err)
        }
        if reply == true {
                http.Redirect(w, r, "/view/success/", http.StatusFound)
        }else {
                http.Redirect(w, r, "/view/fail/", http.StatusFound)
        }
  }  
  func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  	t, _ := template.ParseFiles(tmpl+".html")
  	t.Execute(w,p)
  }
  func check(w http.ResponseWriter, r *http.Request)  {
        client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
        if err != nil {
                log.Fatal("dialing:", err)
        }       
        
        var args = "gggg"
        var reply []byte
        err = client.Call("Req.Get",args, &reply)
        if err != nil {
                log.Fatal("arith error:", err)
        }       
        w.Write(reply)
  }     


  func main() {
  	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/view/success/", succHandler)
	http.HandleFunc("/view/fail/", failHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/save/", saveHandler)
        http.HandleFunc("/login/", loginHandler)
        http.HandleFunc("/login/lg", login)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/post/po", po)
	http.HandleFunc("/check/",check)
  	http.ListenAndServe(":8080", nil)
  }
