.PHONY: build npmi production squad dev

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -gcflags="all=-N -l" -tags lambda.norpc -o bootstrap v1/main.go

npmi:
	npm ci

production: build npmi
	node_modules/.bin/serverless --stage production deploy

squad: build npmi
	node_modules/.bin/serverless --stage squad deploy

dev: build npmi
	node_modules/.bin/serverless --stage dev deploy
