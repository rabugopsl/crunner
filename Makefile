build:
	@go build -o crunner

test:
	@make build
	./crunner --input data.json --debug=true sh echoer.sh
