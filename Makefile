serve: main.go
	go run main.go

clean:
	rm -rf ./**/*.{o,exe}

install:
	go install