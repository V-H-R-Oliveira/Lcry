all:
	go build -ldflags="-s -w" -o winner.lcry app && ../../upx-3.96-amd64_linux/upx --brute winner.lcry
clean:
	rm winner.lcry
