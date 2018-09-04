build: bin/ycp

test:
	go test ./...

bin/ycp:
	dep ensure
	go build -o bin/ycp

run: bin/ycp
	./bin/ycp 

clean:
	rm -rf bin/ycp

.PHONY: build test run clean
