nix:
	go build -o $(GOPATH)/bin/circuit cmd/circuit/main.go

clean:
	rm $(GOPATH)/bin/circuit
