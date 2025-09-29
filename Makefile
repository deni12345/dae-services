MAIN_ENTRY=main

run:
	go run ./cmd/$(MAIN_ENTRY).go

tidy:
	go mod tidy && go mod vendor