BINARY_NAME=tjmsync
PACKAGE=github.com/yezersky/tjmsync
PACKAGE_PATH=/go/src/$(PACKAGE)
IMAGE_NAME_TAG=yezersky/tjmsync

install_dependency:
	glide install

compile_with_docker: install_dependency 
	docker run -v `pwd`:$(PACKAGE_PATH) --rm yezersky/golang-builder -o $(PACKAGE_PATH)/bin/$(BINARY_NAME) $(PACKAGE)

build: compile_with_docker
	docker build -t $(IMAGE_NAME_TAG) .
