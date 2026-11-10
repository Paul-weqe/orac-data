Z3_DIR ?= /home/waswa/projects/golang/z3

frr-grpc:
	protoc --go_out=. --go-grpc_out=. proto/frr/frr-northbound.proto

# Biuld the explore binary
build-grpc:
	go build -o bin/explore cmd/explore/main.go

build-z3:
	CGO_CFLAGS="-I$(Z3_DIR)/src/api" \
	CGO_LDFLAGS="$(Z3_DIR)/build/libz3.so -lpthread" \
	go build -ldflags="-r $(Z3_DIR)/build" -o bin/z3 cmd/z3/main.go

clean:
	rm -rf bin
