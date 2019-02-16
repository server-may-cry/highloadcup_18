.DEFAULT_GOAL := run

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test ./... -bench=. -benchmem

.PHONY: build-server
build-server: test
	GOOS=linux \
	CGO_ENABLED=0 \
	go build -o server

.PHONY: publish
publish: build-server
	docker build -t stor.highloadcup.ru/accounts/happy_dolphin .
	docker push stor.highloadcup.ru/accounts/happy_dolphin

.PHONY: run
run: build-server
	docker build -t stor.highloadcup.ru/accounts/happy_dolphin2 . -f Dockerfile.test.Dockerfile
	docker run -p 8081:80 --rm --memory="2g" --memory-swap="2g" -t stor.highloadcup.ru/accounts/happy_dolphin2
