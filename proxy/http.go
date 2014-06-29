package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var listenHTTP = "0.0.0.0:8881"

func mindHTTP() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}.json", handleStream)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/g/{group:[0-9a-zA-Z-]+}.json", handleGroup)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/m/{media:[0-9a-zA-Z-]+}.json", handleImage)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/m/{media:[0-9a-zA-Z-]+}/{version}.json", handleVersion)
	r.PathPrefix("/").HandlerFunc(handleNotFound)
	log.Fatal(http.ListenAndServe(listenHTTP, r))
}

func handleError(w http.ResponseWriter, r *http.Request, err error, code ...int) {
	if len(code) == 0 {
		code = []int{http.StatusInternalServerError}
	}
	log.Printf("\"%s\" \"%s %s\" %d \"%s\"", r.RemoteAddr, r.Method, r.URL.RequestURI(), code[0], err.Error())
	w.WriteHeader(code[0])
	fmt.Fprint(w, err.Error())
}

func logAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("\"%s\" \"%s %s\" 200 \"OK\"", r.RemoteAddr, r.Method, r.URL.RequestURI())
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	handleError(w, r, fmt.Errorf("not found"), 404)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	logAccess(w, r)
	fmt.Fprint(w,
		`<html>
		<head>
		</head>
		<body>
<pre>
GET /
	Human readable API instructions

GET /{stream}.json
	Fetch JSON data about the stream including group identifiers (which happen to 
	be GUIDs)
	{
		"Name":"the stream name",
		"FirstName":"the stream author firstname",
		"LastName":"the stream author lastname",
		"Groups": ["first group","second group","etc..."]
	}

GET /{stream}/g/{group}.json
	Fetch JSON data about a specific group by its GUID including photo identifiers
	(also GUIDs)
	{
		"Created":"media group creation time",
		"Guid":"media group guid",
		"Caption":"media group caption",
		"FullName":"submitter full name",
		"FirstName":"submitter last name",
		"LastName":"submitter last name",
		"Media":["first item","second item","etc"]
	}

GET /{stream}/m/{item}.json
	Fetch JSON data about a specific media item by its GUID
	{
		"GUID":"media guid",
		"Type":"media type",
		"BatchDateCreated":"media group creation time",
		"BatchGuid":"media group guid",
		"Caption":"media group caption",
		"FullName":"submitter full name",
		"FirstName":"submitter first name",
		"LastName":"submitter last name",
		"Created":"media creation time",
		"Derivatives":{
			"derivative name":{
				"Checksum":"derivative checksum",
				"Size": size in bytes,
				"Height": height in pixels,
				"State":"state (available)",
				"Width": width in pixels
			},
			"another derivative": {...},
			"etc": {},
		}
	}

GET /{stream}/m/{media}/{derivative}.json
	Fetch a JSON list of URLs for a particular derivative.  These signed URLs are
	only valid for a window of time, so fetch them when you are about to use them
	["first url","second url","etc"]

</pre>
		</body>
		</html>`)
}
