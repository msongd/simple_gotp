.PHONY: 

bin/simple_gotp: *.go
	go build -v -o bin/simple_gotp
	
bin/simple_gotp.freebsd: *.go
	env GOOS="freebsd" go build -v -o bin/simple_gotp.freebsd

