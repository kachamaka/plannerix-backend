.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/register ./handlers/register
	env GOOS=linux go build -ldflags="-s -w" -o bin/login ./handlers/login

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
db: 
	dynamodb.sh
