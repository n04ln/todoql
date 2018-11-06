build:
	GOOS=linux GOARCH=amd64 go build -o build/todoql server/server.go

clean:
	rm -rf build

gql:
	mv resolver.go resolver.go_bak
	gqlgen -v

dep:
	dep ensure

dbuild: build
	docker build -t todoql:latest ./

mbuild:
	docker build -f Dockerfile-for-mysql -t todoql_mysql:latest ./
