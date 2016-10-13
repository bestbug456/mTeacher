# Go parameters
GOCMD=go
GOINSTALL=go build
GOARM = env GOOS=linux GOARCH=arm go build
GOLINUX86 = env GOOS=linux GOARCH=386 go build
GOLINUX64 = env GOOS=linux GOARCH=amd64 go build
GODARWIN32 = env GOOS=darwin GOARCH=386 go build
GODARWIN64 = env GOOS=darwin GOARCH=amd64 go build
GOWINDOWS32 = env GOOS=windows GOARCH=386 go build
GOWINDOWS64 = env GOOS=windows GOARCH=amd64 go build
# Avoid problem with the gopath adding
# temporarely ourself to it
GOPATH := $(CURDIR)

# Create a stand-alone repository
# in order to avoid to download
# untested 3° party dependences

GOPATH := $(CURDIR)/src/_vendor:$(GOPATH)

GO_BUILD_ENV := GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH)
#DOCKER_BUILD=$(shell pwd)/.docker_build
#DOCKER_CMD=$(DOCKER_BUILD)/mTeacher

#docker: $(DOCKER_CMD)
#	docker build -t mdota .

#$(DOCKER_CMD): clean
#	mkdir -p $(DOCKER_BUILD)
#	cd ./src && $(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

#clean:
#	rm -rf $(DOCKER_BUILD)

all:
	cd ./src && $(GOINSTALL)
	cd ./src && mv src mTeacher

deploy:
	cd src && $(GOARM) && mv src mTeacher
	mkdir linux_arm && mv src/mTeacher linux_arm
	cd src && $(GOLINUX86) && mv src mTeacher
	mkdir linux_386 && mv src/mTeacher linux_386
	cd src && $(GOLINUX64) && mv src mTeacher
	mkdir linux_amd64 && mv src/mTeacher linux_amd64
	cd src && $(GODARWIN32) && mv src mTeacher
	mkdir darwin_386 && mv src/mTeacher darwin_386
	cd src && $(GODARWIN64) && mv src mTeacher
	mkdir darwin_amd64 && mv src/mTeacher darwin_amd64
	cd src && $(GOWINDOWS32) && mv src.exe mTeacher.exe
	mkdir windows_386 && mv src/mTeacher.exe windows_386
	cd src && $(GOWINDOWS64) && mv src.exe mTeacher.exe
	mkdir windows_amd64 && mv src/mTeacher.exe windows_amd64
	zip -r windows_amd64 windows_amd64
	zip -r windows_386 windows_386
	zip -r darwin_amd64 darwin_amd64
	zip -r darwin_386 darwin_386
	zip -r linux_amd64 linux_amd64
	zip -r linux_386 linux_386
	zip -r linux_arm linux_arm
	rm -rf windows_amd64
	rm -rf windows_386
	rm -rf darwin_amd64
	rm -rf darwin_386
	rm -rf linux_amd64
	rm -rf linux_386
	rm -rf linux_arm