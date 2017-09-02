
setup: ## Install all the build and lint dependencies
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/...
	go get -u github.com/pierrre/gotestcover
	go get -u golang.org/x/tools/cmd/cover
	dep ensure
	gometalinter --install --update

build: ## Build a beta version
	go build -o stowage ./cmd/stowage/

install: ## Install to $GOPATH/src
	go install ./cmd/...

container: Dockerfile stowage ## Create a containerized version of stowage
	docker build -t ealexhudson/stowage .

# refs http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show the different targets available
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
