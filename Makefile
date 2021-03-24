.PHONY: yaml-merge docker-image

BINDIR := ${HOME}/.local/bin

yaml-merge:
	go build -o yaml-merge main.go

install: yaml-merge
	mkdir -p ${BINDIR} >/dev/null
	cp yaml-merge ${BINDIR}

docker-image:
	docker build -t rogstad/yaml-merge .
