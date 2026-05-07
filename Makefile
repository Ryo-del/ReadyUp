PROTO_DIR=api/proto

proto:
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking
