frr-grpc:
	protoc --go_out=. --go-grpc_out=. proto/frr/frr-northbound.proto

# Biuld the explore binary
b-explore:
	go build -o bin/explore cmd/explore/main.go

