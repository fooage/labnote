Use the web page to achieve access to the laboratory log of all platforms, and I usually want to make a log similar to a floating diary. It is a diary way of recording one's own life and emotions, usually exhibiting personal emptiness and unreal life, with pessimistic emotion as the keynote.

## Deploy

Not many things are used, so it is fairly easy to deploy, mainly divided into two parts.

### Config mongodb

Install mongodb on different systems, and you can often customize the database by starting with the configuration file.

```bash
# Please check the official document for specific configuration.
mongod --config ./conf/mongodb.conf
```

- Tips: You can use the `systemctl status mongod` to check the status of mongodb.

### Install server

1. Use commands `go build main.go` to compile executable file.

2. Add html files in the views folder.

## Running

You can change some code in this project to achieve the IP address and port change.

```go
// In ./main.go can change the web server's conncetion info.
Run("127.0.0.1:8090")
// In ./data/data.go can change the mongodb's connection info.
ApplyURI("mongodb://127.0.0.1:27017")
```

If you are not logged in, you will be redirected to the login interface. After successful login, you can write and submit note in the text area of the homepage.

## To Fix

- error redirect caused by 301 cache

Browsers often cache the 301 status, which leads to incorrect redirection after a successful login.

> Solution: `F12` >> `Network` >> `Disable cache`

## To Do

There is lack of a register module in this website, At present, only allowed the administrator add users who is Authorized. But this is also a safer way, because the administrator is a trusted insider!
