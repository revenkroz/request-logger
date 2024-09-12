dev-run:
	cd frontend && yarn build
	go run . -target https://jsonplaceholder.typicode.com

build:
	cd frontend && yarn build
	go build -o requestlogger -ldflags "-w -s" .
