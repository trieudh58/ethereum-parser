# !/bin/make -f

all:
	go build -o build/parserd cmd/parserd/*.go
