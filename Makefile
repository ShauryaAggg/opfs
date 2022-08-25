gen:
	rm -rf pb/pb/*.go
	protoc \
		--proto_path=proto \
		proto/*.proto \
		--go_out=plugins=grpc:pb