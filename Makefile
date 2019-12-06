.PHONY: build
build:
	rm -Rvf test && statik -src=./template && go build .

.PHONY: install
install: build
	go install