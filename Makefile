confmaked: deps
	go build ./cmd/confmaked

deps:
	dep ensure
