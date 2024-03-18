protoc:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/rpc_service_auth.proto

protoc2:
    protoc --go_out=plugins=grpc:. rpc_error_codes.proto

.PHONE: protoc