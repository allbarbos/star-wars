test:
	go test ./... -cover -coverprofile=cover.out

cov: test
	go tool cover -html=cover.out

fmt:
	go fmt ./...
	
api:
	cd api/cmd && go run main.go

importer:
	cd importer/cmd && go run main.go
	