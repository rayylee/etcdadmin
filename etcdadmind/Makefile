Q ?= @

build:
	go build

install:
	go install

uninstall:
	$Q command -v etcdadmind && which etcdadmind | xargs rm -f || :
