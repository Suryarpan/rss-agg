all: run

build:
	go build -o bin/rss-agg

run: build
	./bin/rss-agg

install:
	go install

clean:
	rm -f bin/rss-agg
