SHELL := bash
loc=$(echo ./checkpoints/$(date '+%d-%m-%Y')/)

build:
	cd src && GOOS="linux" GOARCH="amd64" go build -o ../scrape
checkpoint:
	mkdir $(value loc) && cp *json $(value loc)