module github.com/dairaga/webkit

go 1.12

require (
	github.com/dairaga/config v0.0.0-20190611140302-aa0d22deb6c6
	github.com/dairaga/log v0.0.0-20190607012508-d146c6d13bb2
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/schema v1.0.2
	github.com/gorilla/securecookie v1.1.1
)

replace (
	github.com/dairaga/config => ../config
	github.com/dairaga/log => ../log
)
