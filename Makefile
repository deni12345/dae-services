MAIN_ENTRY=main

run:
	go run ./cmd/$(MAIN_ENTRY).go

tidy:
	go mod tidy && go mod vendor

.PHONY: gen
gen:
	@mkdir -p proto/gen
	@protoc \
		-I proto \
		--go_out=proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
		proto/*.proto