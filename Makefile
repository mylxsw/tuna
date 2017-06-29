run:build
	./bin/tuna -conf=./example.toml

build:
	go build -o bin/tuna 

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/tuna-linux 

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/tuna.exe 
