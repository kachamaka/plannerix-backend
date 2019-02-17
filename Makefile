.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/register ./handlers/register
	env GOOS=linux go build -ldflags="-s -w" -o bin/login ./handlers/login
	env GOOS=linux go build -ldflags="-s -w" -o bin/updateSchedule ./handlers/updateSchedule
	env GOOS=linux go build -ldflags="-s -w" -o bin/inputGrade ./handlers/inputGrade
	env GOOS=linux go build -ldflags="-s -w" -o bin/getYearGrades ./handlers/getYearGrades
	env GOOS=linux go build -ldflags="-s -w" -o bin/getWeeklyGrades ./handlers/getWeeklyGrades
	env GOOS=linux go build -ldflags="-s -w" -o bin/getWeeklyEvents ./handlers/getWeeklyEvents
	env GOOS=linux go build -ldflags="-s -w" -o bin/getSubjects ./handlers/getSubjects
	env GOOS=linux go build -ldflags="-s -w" -o bin/getSchedule ./handlers/getSchedule
	env GOOS=linux go build -ldflags="-s -w" -o bin/getNextPeriod ./handlers/getNextPeriod
	env GOOS=linux go build -ldflags="-s -w" -o bin/getAllGrades ./handlers/getAllGrades
	env GOOS=linux go build -ldflags="-s -w" -o bin/getAllEvents ./handlers/getAllEvents
	env GOOS=linux go build -ldflags="-s -w" -o bin/editEvent ./handlers/editEvent
	env GOOS=linux go build -ldflags="-s -w" -o bin/deleteGrade ./handlers/deleteGrade
	env GOOS=linux go build -ldflags="-s -w" -o bin/deleteEvent ./handlers/deleteEvent
	env GOOS=linux go build -ldflags="-s -w" -o bin/createEvent ./handlers/createEvent
	env GOOS=linux go build -ldflags="-s -w" -o bin/changePassword ./handlers/changePassword
	env GOOS=linux go build -ldflags="-s -w" -o bin/changeEmail ./handlers/changeEmail

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
db: 
	dynamodb.sh
