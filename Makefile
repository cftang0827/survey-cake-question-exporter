win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
	zip exporter-win.zip main.exe

mac: 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
	zip exporter-mac.zip main

