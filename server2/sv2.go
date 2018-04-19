  package main
  
  import (
  	"net/http"
        "log"
	"net/rpc"
  )

  const lenPath = len("/view/")
   
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
  	http.HandleFunc("/view/", reqHandler)

  	http.ListenAndServe(":8080", nil)
  }
