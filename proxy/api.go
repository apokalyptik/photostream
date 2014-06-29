package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
			}
			break
		}
	}
	handleError(w, r, fmt.Errorf("not found"), 404)
}
