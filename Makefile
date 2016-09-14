.PHONY: clean

all:
	go build -o bin/totp-cli .

clean:
	rm bin/totp-cli
