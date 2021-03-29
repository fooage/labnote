Use the web page to achieve access to the laboratory log of all platforms, and I usually want to make a log similar to a floating diary. It is a diary way of recording one's own life and emotions, usually exhibiting personal emptiness and unreal life, with pessimistic emotion as the keynote.

## Deploy

Not many things are used, so it is fairly easy to deploy, mainly divided into two parts.

### Config database

1. MongoDB

Install mongodb on different systems, and you can often customize the database by starting with the configuration file.

```bash
# Please check the official document for specific configuration.
mongod --config ./conf/mongodb.conf
```

- Tips: You can use the `systemctl status mongod` to check the status of mongodb.

2. MySQL

This part is being prepared, so stay tuned.

### Install server

1. Use commands `go build main.go` to compile executable file.

2. Add html files in the views folder.

## Running

You can change some code in this project to achieve the IP address and port change, **remember to change your database connection settings**.

```go
// In ./main.go can change the web server's conncetion info.
Run("127.0.0.1:8090")
// And in ./handler/handler.go the cookies should be set to the domain of server.
c.SetCookie("auth", "true", CookieExpireDuration, "/", "127.0.0.1", false, true)
```

If you are not logged in, you will be redirected to the login interface. After successful login, you can write and submit note in the text area of the homepage.

## To Do

Recently, the database and web server have been decoupled, and support for more types of databases will be added.
