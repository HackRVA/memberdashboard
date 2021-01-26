package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/mdns"
)

// ACLCache - store the acl in memory so we can lookup values
var ACLCache []string

func lookupResource() {
	// Make a channel for results and start listening

	qp := mdns.DefaultParams("My awesome service")
	err := mdns.Query(qp)
	if err != nil {
		println(err)
	}
}

func getACL(w http.ResponseWriter, req *http.Request) {
	var ACLReq ACLRequest

	err := json.NewDecoder(req.Body).Decode(&ACLReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(ACLCache)
	w.Write(j)
}

func main() {
	r := mux.NewRouter()

	lookupResource()

	// serve up a frontend that we can test rfid values on
	r.HandleFunc("/gui", serveFiles)
	r.HandleFunc("/get-acl", getACL)
	r.HandleFunc("/", getACL)
	// have an enpoint that accepts acls
	r.HandleFunc("/update", updateHandler)
	// and endpoint to check to see if an rfid value exists
	r.HandleFunc("/lookup", lookupHandler)

	http.Handle("/", r)

	log.Print("Server listening on http://localhost:3001/gui")
	log.Fatal(http.ListenAndServe("0.0.0.0:3001", nil))

}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "./index.html"
	}

	if p == "./gui" {
		p = "./index.html"
	}
	http.ServeFile(w, r, p)
}

// ACLResponse -
// we respond with a hash of the current ACL
type ACLResponse struct {
	Hash string `json:"hash"`
}

// ACLRequest is the json body we expect to receive
//
// {
// 	"acl": [ 2755459513, 848615840 ]
// }
type ACLRequest struct {
	ACL []string `json:"acl"`
}

// ACLHandler takes in the ACL and stores it in a cache
// on the "resource" device this will probably be persisted on the device
func updateHandler(w http.ResponseWriter, req *http.Request) {
	var ACLReq ACLRequest

	err := json.NewDecoder(req.Body).Decode(&ACLReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ACLCache = ACLReq.ACL

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(ACLCache)
	w.Write(j)
}

type lookupRequest struct {
	RFID string `json:"rfid"`
}

type lookupResponse struct {
	Found bool `json:"found"`
}

// lookupHandler - did we find the RFID value?
// note: the actual resources wont' have this endpoint, this is just for easy testing
func lookupHandler(w http.ResponseWriter, req *http.Request) {
	var ack lookupResponse
	var rfid lookupRequest

	err := json.NewDecoder(req.Body).Decode(&rfid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ack.Found = false

	for _, v := range ACLCache {
		if v == rfid.RFID {
			ack.Found = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(ack)
	w.Write(j)
}
