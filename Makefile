version = $(shell cat VERSION)
release_folder=release
linux_build_folder=$(release_folder)/linux-$(arch)-unpacked

executable:
	go build -o garage main.go
build-linux:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o $(linux_build_folder)/garage main.go
	cp config.json $(linux_build_folder)/config.json
package-linux: build-linux
	mkdir -p release
	nfpm pkg --config=nfpm-linux-$(arch).yaml --target release/garage-opener-$(version)-$(arch).deb