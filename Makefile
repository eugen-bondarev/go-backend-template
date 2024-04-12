test:
	@go test -v ./internal/svc/*.go
	@go test -v ./internal/repo/*.go
	@go test -v ./internal/util/*.go