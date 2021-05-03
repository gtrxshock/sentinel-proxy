.PHONY: build run clear
build:
	rm -f build/sentinel-proxy \
	&& GOOS=linux GOARCH=amd64 go build -o build/sentinel-proxy cmd/app/main.go \
	&& ls -l build/sentinel-proxy \
	&& file build/sentinel-proxy
run:
	cp ./build/sentinel-proxy . && ./sentinel-proxy
clear:
	rm -f build/sentinel-proxy
