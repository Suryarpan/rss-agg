all: run

build:
	go build -o bin/rss-agg main.go

run: build
	./bin/rss-agg

install:
	go install

clean:
	rm -f bin/rss-agg
