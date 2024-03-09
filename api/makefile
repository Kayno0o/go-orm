go:
ifeq ($(ENV),prod)
	go build -o ./dist/api ./main.go
	./dist/api
else ifeq ($(ENV),dev)
	air -c .air.toml
else
	@echo "Invalid environment specified. Use 'env=prod' or 'env=dev'."
	@exit 1
endif

build-pug:
	rm -rf src/jade/*
	jade -writer -pkg=jade -d=src/jade /app/client/pages/*.pug

pug:
ifeq ($(ENV),prod)
	make build-pug
else ifeq ($(ENV),dev)
	./pug_watch.sh
else
	@echo "Invalid environment specified. Use 'env=prod' or 'env=dev'."
	@exit 1
endif
