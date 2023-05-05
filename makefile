test:
	go test -v --cover ./...

protoc: 
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

start:
	go run main.go

grpc_mock:
	mockgen -package mockpb -destination pb/mock/client.go github.com/machearn/galaxy_controller/pb GalaxyClient

.PHONY:
	test protoc start grpc_mock