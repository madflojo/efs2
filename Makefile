# Efs2 Makefile is used to drive the build and installation of Efs2
# this is meant to be used with a local copy of code repository.

tests:
	@echo "Launching Tests in Docker Compose"
	mkdir -p ./testdata
	test -f ./testdata/testkey || ssh-keygen -b 2048 -t rsa -f ./testdata/testkey -q -N ""
	test -f ./testdata/wrongkey || ssh-keygen -b 2048 -t rsa -f ./testdata/wrongkey -q -N ""
	test -f ./testdata/invalidkey || ssh-keygen -b 2048 -t rsa -f ./testdata/invalidkey -q -N "" && sed '5d' ./testdata/invalidkey > ./testdata/tmp.key && mv ./testdata/tmp.key ./testdata/invalidkey
	test -f ./testdata/testkey-passphrase || ssh-keygen -b 2048 -t rsa -m PEM -f ./testdata/testkey-passphrase -q -N "testing"
	docker-compose -f dev-compose.yml up --build tests

clean:
	@echo "Cleaning up build junk"
	-docker-compose -f dev-compose.yml down
	-rm -rf ./testdata

build:
	@echo "Building from source"
	go build

install:
	@echo "Installing from source"
	go install
