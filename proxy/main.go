package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var listenHttp = "0.0.0.0:8881"

func init() {
	flag.StringVar(&listenHttp, "http", listenHttp, "http address and port number to listen on")
	flag.DurationVar(&cacheDuration, "cache", cacheDuration, "keep items in the local cache for this long")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	logAccess(w, r)
	fmt.Fprint(w,
		`<html>
		<head>
		</head>
		<body>
<pre>
GET /{stream}.json
	{
		"Name":"the stream name",
		"FirstName":"the stream author firstname",
		"LastName":"the stream author lastname",
		"Groups": ["first group","second group","etc..."]
	}

GET /{stream}/{group}.json
	{
		"Created":"media group creation time",
		"Guid":"media group guid",
		"Caption":"media group caption",
		"FullName":"submitter full name",
		"FirstName":"submitter last name",
		"LastName":"submitter last name",
		"Media":["first item","second item","etc"]
	}

GET /{stream}/{group}/{item}.json
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

GET /{stream}/{group}/{media}/{derivative}
	["first url","second url","etc"]

</pre>
		</body>
		</html>`)
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

func handleStream(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   vars["stream"],
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		if err.Error() == "unexpected response: 404 Not Found" {
			handleNotFound(w, r)
		} else {
			handleError(w, r, err)
		}
		return
	}
	var rsp = struct {
		Name      string
		FirstName string
		LastName  string
		Groups    []string
	}{
		Name:      stream.Name,
		FirstName: stream.FirstName,
		LastName:  stream.LastName,
	}
	for k, v := range stream.Media {
		if len(stream.Media) != (k+1) && v.BatchGuid == stream.Media[(k+1)].BatchGuid {
			continue
		}
		rsp.Groups = append(rsp.Groups, v.BatchGuid)
	}
	if bytes, err := json.Marshal(rsp); err != nil {
		handleError(w, r, err)
	} else {
		logAccess(w, r)
		w.Write(bytes)
	}
}

func handleGroup(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   vars["stream"],
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		handleError(w, r, err)
		return
	}
	var rsp = struct {
		Created   string
		Guid      string
		Caption   string
		FullName  string
		FirstName string
		LastName  string
		Media     []string
	}{}
	for k, v := range stream.Media {
		if v.BatchGuid != vars["group"] {
			continue
		}
		rsp.Media = append(rsp.Media, v.GUID)
		if len(stream.Media) != (k+1) && stream.Media[(k+1)].BatchGuid == v.BatchGuid {
			continue
		}
		rsp.Created = v.BatchDateCreated
		rsp.Guid = v.BatchGuid
		rsp.Caption = v.Caption
		rsp.FullName = v.FullName
		rsp.FirstName = v.FirstName
		rsp.LastName = v.LastName
	}
	if rsp.Media == nil {
		handleNotFound(w, r)
		return
	}
	if bytes, err := json.Marshal(rsp); err != nil {
		handleError(w, r, err)
	} else {
		logAccess(w, r)
		w.Write(bytes)
	}
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   vars["stream"],
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		handleError(w, r, err)
		return
	}
	for _, v := range stream.Media {
		if v.GUID == vars["media"] {
			if bytes, err := json.Marshal(v); err != nil {
				handleError(w, r, err)
			} else {
				logAccess(w, r)
				w.Write(bytes)
			}
			return
		}
	}
	handleError(w, r, fmt.Errorf("not found"), 404)
}
func handleVersion(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   vars["stream"],
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		handleError(w, r, err)
		return
	}
	for _, v := range stream.Media {
		if v.GUID == vars["media"] {
			if d, ok := v.Derivatives[vars["version"]]; ok {
				if u, e := d.GetURLs(); e != nil {
					handleError(w, r, e)
				} else {
					if bytes, err := json.Marshal(u); err != nil {
						handleError(w, r, err)
					} else {
						logAccess(w, r)
						w.Write(bytes)
					}
				}
				return
			} else {
				break
			}
		}
	}
	handleError(w, r, fmt.Errorf("not found"), 404)
}

func main() {
	flag.Parse()
	go mindEngine()
	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}.json", handleStream)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/{group:[0-9a-zA-Z-]+}.json", handleGroup)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/{group:[0-9a-zA-Z-]+}/{media:[0-9a-zA-Z-]+}.json", handleImage)
	r.HandleFunc("/{stream:[0-9a-zA-Z]+}/{group:[0-9a-zA-Z-]+}/{media:[0-9a-zA-Z-]+}/{version}", handleVersion)
	r.PathPrefix("/").HandlerFunc(handleNotFound)
	log.Fatal(http.ListenAndServe(listenHttp, r))
}
