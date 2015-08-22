.PHONY: all

PROJECT = craftctl
VERSION = 0.1.0

all: build

clean:
	$(RM) -rf dist/

build:
	scripts/build.sh

release:
	git push origin --tags
	python scripts/release.py $(PROJECT) $(VERSION)
