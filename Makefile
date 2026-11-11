frr-grpc:
	protoc --go_out=. --go-grpc_out=. proto/frr/frr-northbound.proto

# Biuld the explore binary
build:
	go build -o bin/explore cmd/explore/main.go

clean:
	rm -rf bin
