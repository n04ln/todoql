# todoql

## Required
- <a href="https://github.com/golang/dep">github.com/golang/dep</a>

## Run server
``` sh
$ dep ensure && make clean && make build && make dbuild && make mbuild && docker-compose down && docker-compose up
```

go to <a href='http://localhost:8080'>http://localhost:8080</a>

using pre-data (in `provisioning.sql`), can get data using data-loader

### query
```
query findTodo{
  getUser(id:"qwer") {
    name
    todos{
      id
      text
      user{
        name
        todos{
          text
        }
      }
      done
    }
  }
}
```
