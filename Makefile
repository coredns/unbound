VERSION:=0.0.6
TAG:=v$(VERSION)

all:
	@echo Use the 'release' target to start a release $(VERSION)

.PHONY: release
release: commit push
	@echo Released $(VERSION)

.PHONY: commit
commit:
	@echo Committing release $(VERSION)
	git commit -am"Release $(VERSION)"
	git tag $(TAG)

.PHONY: push
push:
	@echo Pushing release $(VERSION) to master
	git push --tags
	git push
