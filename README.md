# dag-go
> Job scheduler written with Go

## Feature
- Generate dag based on XML config and execute

## How To Build & Run
```sh
$ make build
$ ./dag-go -xmlFilePath=[xmlFilePath]
```

## Example
```sh
$ make build
$ ./dag-go -xmlFilePath=./xml/test03.xml
```

## TODO
- Execute shell command in Docker container
- Webserver for monitoring tasks