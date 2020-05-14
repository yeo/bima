tag ?= v0.1.0

build:
	go build -o bima
	fyne package -os darwin -icon ../../assets/icon.png
	rm -rf ~/Applications/bima.app
	mv bima.app ~/Applications

release:
	github-release release \
    --user yeo \
    --repo bima \
    --tag $(tag) \
    --name "$(tag)" \
    --description "Laputa"

	github-release upload \
    --user yeo \
    --repo bima \
    --tag $(tag) \
    --name "bima-osx.tar.gz" \
    --file "bima.tar.gz"

