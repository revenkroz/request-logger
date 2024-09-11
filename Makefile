dev-run:
	cd frontend && yarn build
	go run . -url https://example.com
