.PHONY: 

bin/simple_gotp: *.go frontend/js/* frontend/css/* frontend/img/* frontend/*
	go build -v -o bin/simple_gotp
	
bin/otp_cli: cli/*.go
	go build -v -o bin/otp_cli cli/*.go
