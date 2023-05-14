clean:
	rm -f aos

build: clean
	go build -o aos cmd/main.go

example: build
	./aos host -c example.yml -o result.json

example-image: build
	./aos image registry.altlinux.org/alt/alt -c example.yml -o result.json