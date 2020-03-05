.PHONY: test dev-tools ci-tools tools commit

test:
	go test ./...

tools: dev-tools ci-tools

dev-tools:
	npm install commitizen@4.0.3 \
		@commitlint/cli@8.3.5 \
		@commitlint/config-conventional@8.3.4

ci-tools:
	npm install semantic-release@16.0.0 \
		@semantic-release/changelog@3.0.6 \
		@semantic-release/commit-analyzer@7.0.0 \
		@semantic-release/git@8.0.0 \
		@semantic-release/github@6.0.0 \
		@semantic-release/release-notes-generator@7.3.5

commit:
	npx git-cz
