frr-grpc:
	protoc --go_out=. --go-grpc_out=. proto/frr/frr-northbound.proto
