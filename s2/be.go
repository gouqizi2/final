  package main

  import (
  	"net/http"
        "log"
	"net/rpc"
	"io/ioutil"
	"net"
    "strings"
  )

var cache = make(map[string]string)
type Req struct{
}

  const lenPath = len("/view/")
  type Page struct {
        Title   string
       	Body    []byte
  }
  type Pass struct {
	Username string
	Page   Page
  }
var pn *paxosNetwork
var as = make([]*acceptor,0)
var cont = make(map[string]Page)
  type Info struct {
	Username string
	Password string
  }
  func (p *Page) save() {
        filename := p.Title + ".txt"
        ioutil.WriteFile(filename, p.Body, 0600)
  }

  func (r *Req) Create(i Info, reply *bool) error {
    p := newProposer(1001, i.Password , pn.agentNetwork(1001), 1, 2, 3)
    go p.run()

    l := newLearner(2001, pn.agentNetwork(2001), 1,2,3)
    value := l.learn()
	cache[i.Username] = value
	*reply = true
	return nil
  } 
  func (r *Req) Login(i Info, reply *bool) error {
    p:= newProposer(1001, cache[i.Username] , pn.agentNetwork(1001), 1, 2, 3)
    go p.run()
    l := newLearner(2001, pn.agentNetwork(2001), 1,2,3)
    value := l.learn()
	if i.Password == value {
		*reply = true
	}else   {
		*reply = false
	}
	return nil
  }
  func (r *Req) Post(p Pass, reply *bool) error {
    p1 := newProposer(1001, p.Page.Title + ":" + (string) (p.Page.Body) , pn.agentNetwork(1001), 1, 2, 3)
    go p1.run()
    l := newLearner(2001, pn.agentNetwork(2001), 1,2,3)
    value := l.learn()
    cont[p.Username] = Page{Title:(strings.Split(value,":")[0]), Body:([]byte)(strings.Split(value,":")[1]) }

    *reply = true
    return nil
  }
  func (r *Req) Get(usern string, reply *[]byte) error {
	var tempres = ""
	for u,v := range cont {
		tempres = tempres + u + ":" + v.Title+" " + string(v.Body) + "\n"
	}
	*reply = []byte(tempres)
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
    pt := newPaxosNetwork(1, 2, 3, 1001, 2001)
    pn = pt
 //   as := make([]*acceptor, 0)
    for i := 1; i <= 3; i++ {
        as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
    }
    for _, a := range as {
        go a.run()
    }
	rpc.Register(new(Req))
	rpc.HandleHTTP()
    	l, e := net.Listen("tcp", ":1234")
    	if e != nil {
        	log.Fatal("listen error:", e)
  	}
        http.Serve(l, nil)
  }
