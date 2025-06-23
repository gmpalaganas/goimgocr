run:
	go run ./cmd/main/main.go $(ARGS)

build:
	go build -o ./build/giocr ./cmd/main/main.go

install:
	go build -o ~/go/bin/giocr ./cmd/main/main.go

clean:
	rm -rf ./build/

uninstall:
	rm ~/go/bin/giocr
