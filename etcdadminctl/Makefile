Q ?= @

build:
	go build

install:
	go install

uninstall:
	$Q command -v etcdadminctl && which etcdadminctl | xargs rm -f || :
