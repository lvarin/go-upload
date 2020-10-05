all: fileupload

fileupload: fileupload.go
	go build fileupload.go
