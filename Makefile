.PHONY: clean build deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o bin/register ./handlers/register
	env GOOS=linux GOARCH=amd64 go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o bin/login ./handlers/login
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/updateSchedule ./handlers/updateSchedule
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/inputGrade ./handlers/inputGrade
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getYearGrades ./handlers/getYearGrades
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getWeeklyGrades ./handlers/getWeeklyGrades
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getWeeklyEvents ./handlers/getWeeklyEvents
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getNextPeriod ./handlers/getNextPeriod
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getAllGrades ./handlers/getAllGrades
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getAllEvents ./handlers/getAllEvents
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/editEvent ./handlers/editEvent
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/deleteGrade ./handlers/deleteGrade
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/deleteEvent ./handlers/deleteEvent
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/createEvent ./handlers/createEvent
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/changePassword ./handlers/changePassword
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/changeEmail ./handlers/changeEmail
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getProfile ./handlers/getProfile
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/updateNotifications ./handlers/updateNotifications
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/registerUnverified ./handlers/registerUnverified
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/checkVerified ./handlers/checkVerified
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/sendEmail ./handlers/sendEmail

	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getSchedule ./handlers/getSchedule
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/putSchedule ./handlers/putSchedule

	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/getSubjects ./handlers/getSubjects
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/updateSubjects ./handlers/updateSubjects
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/createSubjects ./handlers/createSubjects
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/deleteSubjects ./handlers/deleteSubjects

	

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --noDeploy
	./go-serverless
	aws cloudformation deploy --template-file ./.serverless/cloudformation-template-update-stack.json --stack-name plannerix --s3-bucket plannerix-kinghunter58 --force-upload

db-t:
	java -Djava.library.path=/home/trayan/Code/dynamodb/DynamoDBLocal_lib -jar /home/trayan/Code/dynamodb/DynamoDBLocal.jar -dbPath ./resources/ -sharedDb

db-local:
	dynamodb-admin -o