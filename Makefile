build: bin/ycp

test:
	GO111MODULE=on go test ./...

bin/ycp:
	GO111MODULE=on go build -o bin/ycp

run: bin/ycp
	./bin/ycp 

clean:
	rm -rf bin/ycp

.PHONY: build test run clean
