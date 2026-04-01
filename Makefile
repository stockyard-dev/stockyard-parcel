build:
	CGO_ENABLED=0 go build -o parcel ./cmd/parcel/

run: build
	./parcel

test:
	go test ./...

clean:
	rm -f parcel

.PHONY: build run test clean
