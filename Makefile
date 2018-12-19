.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello ./handlers/hello
	env GOOS=linux go build -ldflags="-s -w" -o bin/world ./handlers/world

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
db: 
	dynamodb.sh
