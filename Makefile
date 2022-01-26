start: main.go
	go run main.go

serve: main.go
	air
	#nodemon --exec go run main.go --signal SIGTERM


clean:
	rm -rf ./**/*.{o,exe}

install:
	go install