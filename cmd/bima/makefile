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

clean:
	rm -rf Payload bima.app bima.ipa

ios:
	fyne package -os ios -appID space.yeo.bima -icon assets/icon.png


dev-build:
	codesign -f -s 'Apple Development: Vinh Nguyen (63C8TVC926)' bima.app
	xcrun simctl install booted bima.app

prod-build:
	codesign -f -vv -s 'Apple Distribution: VINH  Nguyen (456HYLV3SG)' --entitlements entitlements.plist Payload/bima.app

ipa:
	# strip
	# without this we will get not suppoer error
	lipo bima.app/main -info
	(lipo bima.app/main -remove i386 -remove x86_64 -output bima.app/strip-main && mv bima.app/strip-main bima.app/main) || true
	# https://apple.stackexchange.com/questions/173471/changing-info-plist-values-via-terminal
	# Without this we willg et ipad error
	#plutil -insert UIRequiresFullScreen -bool YES bima.app/Info.plist || true
	./bump-version.sh patch bima.app/Info.plist
	#cp Info.plist bima.app/
	# Copy icon
	cp assets/*.png bima.app/
	#cp assets/*.png bima.app
	mkdir Payload || true
	#cp -r Assets.xcassets bima.app/
	cp -r Images.xcassets bima.app/
	cp -rf bima.app Payload/
	cp '/Users/vinh/Library/Mobile Documents/com~apple~CloudDocs/bima-ios/Bima_Production_Distribution.mobileprovision' Payload/bima.app/embedded.mobileprovision
	codesign -f -vv -s 'Apple Distribution: VINH  Nguyen (456HYLV3SG)' --entitlements entitlements.plist Payload/bima.app
	zip -r bima.ipa Payload/

icon:
	convert icon.png -resize 120x120 icon-60@2x.png
	convert icon.png -resize 180x180 icon-60@3x.png
	convert icon.png -resize 152x152 icon-76@2x.png
	convert icon.png -resize 76x76 icon-76.png
