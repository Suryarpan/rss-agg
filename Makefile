all: run

build:
	go build -o bin/rss-agg main.go

run: build
	./bin/rss-agg

clean:
	rm -f bin/rss-agg
