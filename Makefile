SOURCEDIR=.
SOURCES :=  $(shell find $(SOURCEDIR) -name '*_test.go')

OUTDIR=bin/
BINARY=data

VERSION=1.0.0
BUILD=`date +"%s"`

LDFLAGS=-ldflags "-X github.com/rithium/version.Version=${VERSION} -X github.com/rithium/version.Build=${BUILD}"

build:
	go get -d -v
	go build -v ${LDFLAGS} -o ${OUTDIR}${BINARY} main.go

xbuild:
	env GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o ${OUTDIR}${BINARY}-linux main.go

test:
	go get -d -v
	go test -tags=integration

install:
	echo "install"
	go install ${LDFLAGS}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install