# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test

test: 
	$(GOTEST) -v ./...

# release new version
# make release tag=v1.x.x
release:
	./release.sh $(tag)