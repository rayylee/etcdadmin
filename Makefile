Q ?= @

build: build-pre
	go build
	make -C etcdadminctl build

install: build-pre
	go install
	make -C etcdadminctl install

uninstall:
	$Q command -v etcdadmin && which etcdadmin | xargs rm -f || :
	make -C etcdadminctl uninstall

all: clean install

contrib: all
	bash contrib/install
	bash build contrib

uncontrib: uninstall
	bash contrib/install uninstall
	bash build uncontrib

build-pre:
	bash build build-pre

clean:
	$Q rm -f etcdadminctl/etcdadminctl || :
	$Q rm -f etcdadmin || :

.PHONY: contrib build
