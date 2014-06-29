Photostream Proxy
=================
A web based proxy to make consuming public photostreams super easy

# server usage

```
Usage of proxy:
  -cache=5m0s: keep items in the local cache for this long
  -gid=0: set GID (0 disables)
  -http="0.0.0.0:8881": http address and port number to listen on
  -uid=0: set UID (0 disables)
```

# client usage

#### GET /
	Human readable API instructions

#### GET /{stream}.json
	Fetch JSON data about the stream including group identifiers (which happen to 
	be GUIDs)
```json
	{
		"Name":"the stream name",
		"FirstName":"the stream author firstname",
		"LastName":"the stream author lastname",
		"Groups": ["first group","second group","etc..."]
	}
```

#### GET /{stream}/g/{group}.json
	Fetch JSON data about a specific group by its GUID including photo identifiers
	(also GUIDs)
```json
	{
		"Created":"media group creation time",
		"Guid":"media group guid",
		"Caption":"media group caption",
		"FullName":"submitter full name",
		"FirstName":"submitter last name",
		"LastName":"submitter last name",
		"Media":["first item","second item","etc"]
	}
```

#### GET /{stream}/m/{item}.json
	Fetch JSON data about a specific media item by its GUID
```json
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
```

#### GET /{stream}/m/{media}/{derivative}.json
	Fetch a JSON list of URLs for a particular derivative.  These signed URLs are
	only valid for a window of time, so fetch them when you are about to use them
```json
	["first url","second url","etc"]
```
