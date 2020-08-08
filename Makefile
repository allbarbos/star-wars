GOCMD=go
PROJECT=star-wars

test:
	go test ./... -cover -coverprofile=$(PROJECT).coverprofile

cov: test
	go tool cover -html=$(PROJECT).coverprofile

api:
	cd api/cmd && go run main.go

importer:
	cd importer/cmd && go run main.go
