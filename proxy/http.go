package main

import "github.com/gin-gonic/gin"

var listenHTTP = "0.0.0.0:8881"

func mindHTTP() {
	r := gin.Default()
	r.GET("/", handleIndex)
	r.GET("/:stream", handleStream)
	r.GET("/:stream/g/:group", handleGroup)
	r.GET("/:stream/m/:media", handleMedia)
	r.GET("/:stream/m/:media/:version", handleVersion)
	r.Run(listenHTTP)
}

func handleIndex(c *gin.Context) {
	c.String(200,
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
