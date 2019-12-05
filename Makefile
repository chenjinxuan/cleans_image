# build linux
build:
	GOOS=linux GOARCH=amd64 go build -o server main.go
