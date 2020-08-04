GOCMD=go
PROJECT=star-wars

test:
	go test ./... -cover -coverprofile=$(PROJECT).coverprofile
	# go test $(go list ./... | grep -v /repository | grep -v /utils | grep -v /entity/adapter)

cov: test
	go tool cover -html=$(PROJECT).coverprofile
