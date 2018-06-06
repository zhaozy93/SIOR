GOPATH:=$(CURDIR)
export GOPATH

all: build

fmt:
	gofmt -l -w -s src/

dep:fmt
	#go get github.com/shengkehua/xlog4go
	#go get github.com/garyburd/redigo/redis
	#go get code.google.com/p/gcfg
	#go get git.apache.org/thrift.git/lib/go/thrift
	#go get go.intra.xiaojukeji.com/golang/protobuf/proto
	#go get github.com/Shopify/sarama
	#go get go.intra.xiaojukeji.com/apollo/apollo-golang-sdk
	#go get go.intra.xiaojukeji.com/engine/zmq4
	#go get go.intra.xiaojukeji.com/shield-arch/simhelper

build:dep
	go build -o bin/stdriverstatus main
client:
	go build -o bin/client src/client/client.go
	chmod +x bin/client

clean:
	#rm -rfv pkg
	rm -rf bin/stdriverstatus
	rm -rf status
	rm -rf output
output: build
	mkdir -p output/bin
	mkdir -p output/conf
	mkdir -p output/log
	mkdir -p output/status
	mkdir -p output/citydata
	cp -r bin/* output/bin/
	cp -r conf/* output/conf/
	cp -r citydata/* output/citydata
	cp load.sh output/

deploy: 
	rm -rf output/bin
	rm -rf output/conf
	rm -rf output/log
	mkdir -p output
	cp -r bin output/
	cp -r conf output/
	cp -r citydata output/
	cp control.sh output/
	cp loadbroadcastdist.py output/
	cp -R deploy-meta output/
	cp Dockerfile output/

start:
	go run src/main/main.go
