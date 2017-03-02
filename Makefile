BINARY_NAME=tjmsync
PACKAGE=github.com/yezersky/tjmsync/bin/tjmsync
PACKAGE_PATH=/go/src/github.com/yezersky/tjmsync
IMAGE_NAME_TAG=yezersky/tjmsync

install_dependency:
	glide install

compile_with_docker: install_dependency 
	docker run -v `pwd`:$(PACKAGE_PATH) --rm golang:1.7 go build -o $(PACKAGE_PATH)/out/$(BINARY_NAME) $(PACKAGE)

build: compile_with_docker
	docker build -t $(IMAGE_NAME_TAG) .

clean:
	rm out/tjmsync
