BINARY_NAME = proxy

all: build

multi-platform: build build-linux build-windows build-windows_64 build_mac_386 build_mac_64

build:
	go get -t -v
	go build -o $(BINARY_NAME) -v

run:
	go get -t -v
	go build -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

build-linux:
	go get -t -v
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -o $(BINARY_NAME)_linux -v

build-windows:
	go get -t -v
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -installsuffix cgo -o $(BINARY_NAME)_windows_386.exe -v

build-windows_64:
	go get -t -v
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARY_NAME)_windows_64.exe -v

build_mac_386:
	go get -t -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -installsuffix cgo -o $(BINARY_NAME)_mac_386 -v

build_mac_64:
	go get -t -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o $(BINARY_NAME)_mac_64 -v


clean:
	rm $(BINARY_NAME)*
	rm -f ./temp/*