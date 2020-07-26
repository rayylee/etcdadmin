Q ?= @

build: build-pre
	make -C etcdadmind build
	make -C etcdadminctl build

install: build-pre
	make -C etcdadmind install
	make -C etcdadminctl install

uninstall:
	make -C etcdadmind uninstall
	make -C etcdadminctl uninstall

all: clean install

contrib: all
	bash contrib/install
	bash build contrib

uncontrib: uninstall
	bash contrib/install uninstall
	bash build uncontrib

build-pre:
	mkdir -p etcdadmind/pb/etcdadminpb || :
	mkdir -p etcdadminctl/pb/etcdadminpb || :
	bash build build-pre

clean:
	$Q rm -f etcdadminctl/etcdadminctl || :
	$Q rm -f etcdadmind/etcdadmind || :

.PHONY: contrib build
