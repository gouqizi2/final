  package main
  
  import (
  	"net/http"
        "log"
	"net/rpc"
	"io/ioutil"
	"net"
  )

var cache = make(map[string]string)

type Req struct{
}

  const lenPath = len("/view/")
    type Page struct {
        Title   string
       	Body    []byte
  }
  type Info struct {
	Username string
	Password string
  }
  func (p *Page) save() {
        filename := p.Title + ".txt"
        ioutil.WriteFile(filename, p.Body, 0600)
  }

  func (r *Req) Create(i Info, reply *bool) error {
	cache[i.Username] = i.Password
	*reply = true
	return nil
  } 
  func (r *Req) Login(i Info, reply *bool) error {
	if i.Password == cache[i.Username] {
		*reply = true
	}else   {
		*reply = false
	}
	return nil
  }
//  func (r *Req) Create(args Page, reply *int) error {
//	temp := &page{title: args.Title, body:[]byte(args.Body)}
//	
//		args.save()
//		*reply = 1
//		return nil	
//	}

//   func (r *Req) Login(args string, reply *bool) error {
//
//        b, _ := (ioutil.ReadFile("account.txt"))
//        b2 := string(b[:])
//        if b2 == args {
//                *reply = true
//        }else {
//                *reply = false
//        }
//	return nil
//  }
  func reqHandler(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
   	if err != nil {
        	log.Fatal("dialing:", err)
   	}
	err = client.Call("Req.Process", title, &w)
	if err != nil {
        	log.Fatal("arith error:", err)
    	}	
  }

  func main() {
	rpc.Register(new(Req))
	rpc.HandleHTTP()
    	l, e := net.Listen("tcp", ":1234")
    	if e != nil {
        	log.Fatal("listen error:", e)
  	}
        http.Serve(l, nil)
  }
