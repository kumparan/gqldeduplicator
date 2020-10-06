changelog_args=-o CHANGELOG.md -p '^v'

lint:
	golangci-lint run --exclude-use-default=false --enable=golint --enable=goimports --enable=unconvert --enable=unparam --enable=gosec

test:
	go test -v --cover .

changelog:
ifdef version
	$(eval changelog_args=--next-tag $(version) $(changelog_args))
endif
	git-chglog $(changelog_args)

.PHONY: lint test changelog
