all:
	go build -ldflags="-s -w" -o moon.lcry app && ../../upx-3.96-amd64_linux/upx --brute moon.lcry
clean:
	rm moon.lcry
