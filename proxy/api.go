package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleStream(c *gin.Context) {
	var reqstream = c.Params.ByName("stream")
	if reqstream[len(reqstream)-5:] != ".json" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	reqstream = reqstream[:len(reqstream)-5]
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   reqstream,
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		if err.Error() == "unexpected response: 404 Not Found" {
			c.String(http.StatusNotFound, "not found")
		} else {
			c.String(http.StatusInternalServerError, err.Error())
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
	c.JSON(200, rsp)
}

func handleGroup(c *gin.Context) {
	var reqstream = c.Params.ByName("stream")
	var reqgroup = c.Params.ByName("group")
	if reqgroup[len(reqgroup)-5:] != ".json" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	reqgroup = reqgroup[:len(reqgroup)-5]
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   reqstream,
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
		if v.BatchGuid != reqgroup {
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
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.JSON(200, rsp)
}

func handleMedia(c *gin.Context) {
	var reqstream = c.Params.ByName("stream")
	var reqmedia = c.Params.ByName("media")
	if reqmedia[len(reqmedia)-5:] != ".json" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	reqmedia = reqmedia[:len(reqmedia)-5]
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   reqstream,
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range stream.Media {
		if v.GUID == reqmedia {
			c.JSON(200, v)
			return
		}
	}
	c.String(http.StatusNotFound, "not found")
}

func handleVersion(c *gin.Context) {
	var reqstream = c.Params.ByName("stream")
	var reqmedia = c.Params.ByName("media")
	var reqversion = c.Params.ByName("version")
	if reqversion[len(reqversion)-5:] != ".json" {
		c.String(http.StatusNotFound, "not found")
		return
	}
	reqversion = reqversion[:len(reqversion)-5]
	var reply = make(chan *cacheEntry)
	fetch <- request{
		stream:   reqstream,
		response: reply,
	}
	var client = <-reply
	stream, err := client.getStream()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range stream.Media {
		if v.GUID == reqmedia {
			if d, ok := v.Derivatives[reqversion]; ok {
				if u, e := d.GetURLs(); e != nil {
					c.String(http.StatusInternalServerError, e.Error())
				} else {
					c.JSON(200, u)
				}
				return
			}
			break
		}
	}
	c.String(http.StatusNotFound, "not found")
}
