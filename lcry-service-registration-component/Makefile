all:
	go build -ldflags="-s -w" -o loader.lcry app && ../../upx-3.96-amd64_linux/upx --brute loader.lcry
makeTar:
	tar -cvzf - moon.lcry winner.lcry | openssl enc -pbkdf2 -e -aes256 -out lcry.tgz
clean:
	rm loader.lcry
