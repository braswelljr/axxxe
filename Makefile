.PHONY: start
start: main.go
	go run main.go

.PHONY: serve
serve: main.go
	air

.PHONY: nodemon
nodemon: main.go
	nodemon --exec go run main.go --signal SIGTERM

# listen and kill process using port 3030
clearport:
	lsof -i TCP:3030 | grep LISTEN | awk '{print $2}' | xargs kill
	kill -9 $(lsof -i TCP:3030 | grep LISTEN | awk '{print $2}')

startmongodb:
	systemctl start mongod.service

clean:
	rm -rf ./**/*.{o,exe}

install:
	go install

tidy:
	go fmt ./**/*.go
