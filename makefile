test:
	go test -v --cover ./...

protoc: 
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

start:
	go run main.go

.PHONY:
	test protoc start