# hmn

> The development of this project is on pause for now,
> there are some bug not yet fixed and the whole project is insecure and unpolished at the moment

> !!This project is under development and is insecure and unpolished at the moment!!



hmn is hierarchical note taking application like [Trilium Notes](https://github.com/zadam/trilium) but without all the function that I don't need.

I have been using Trilium Notes for a few year but I found myself not using most of the feature it provide, so I would like to make an alternative that only contain the feature that I need. I am hoping that by writing this application in golang I can reduce its memory usage so that I can run other useful selfhost application on my respberry pi

![screenshot of hmn](https://rchan.codeberg.page/hmn/img/hmn-main.png)

## Unimplemented/unpolished Feature
- security is not priority at current stage of development
- error handling
- efficiently send update to server
- the entire front-end could be better



## Repository Structure


Directory | Description
------------ | -------------
[`business/`](business/) | Application tier logic. 
[`cli/`](cli/) | Logic for handling command line argument
[`config/`](config/) | Structure and helper function for loading config file
[`constant/`](constant/) | Constant value
[`repository/`](repository/) | Data tier logic and helper to create new sqlite data store
[`dto/`](dto/) | Dto for transmitting data with web controller or repository
[`helper/`](helper/) | Helper function that used by different part of the application
[`log/`](log/) | Provide a zap log instance for other module to call
[`model/`](model/) | Business model that might have business logic in it
[`public/`](public/) | Front-end code
[`web/`](web/) | Web server, web controller and middleware

## How to run 


### Example config file

`config.json`
```
{
    "storage":{
        "type":"sqlite3",
        "path":"db.sqlite3"
    },
    "server":{
        "use_https":false,
        "tls_key":null,
        "tls_cert":null,
        "listen_on":":8080",
        "allow_origins":[]
    }
}
```

### Run without build

```
go run main.go start
```

### Build and Run
```
go build

./hmn start
```

### Build for Window on Linux
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc  go build
```

### Set password for user

You need to setup the password for user at first run.

When application is not running
```
go run main.go reset -u <username>
```
username is admin by default


## Feature 
- create new note under by clicking `+` after the note title
- drag note onto another note's title to put it under the note
- drag note to space between two note title to move the note to that spot
- cannot move parent to under itself


## External Library
- [echo](https://github.com/labstack/echo)
- [cobra](https://github.com/spf13/cobra)
- [simpleMDE](https://github.com/sparksuite/simplemde-markdown-editor)
- [preact](https://github.com/preactjs/preact)
- [sqlx](https://github.com/jmoiron/sqlx)
