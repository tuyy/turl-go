.PHONY: build clean

build:
	mkdir -p dist
	env GOOS=linux GOARCH=amd64 go build -o dist/turlForCentos ./cmd/turl
	go build -o dist/turl ./cmd/turl

clean:
	rm -rf dist