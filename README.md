# crawl
VERY unfinished game. The idea is to make an open-world terminal-based multiplayer game possible. I'm trying to make it as programmable as possible, but a level editor is beyond the scope right now. To make your own world, take a look at server.go (this is where the actual game content is defined).

To run:
```
go run main.go -server #to run the server

go run main.go -addr <ip>:9000 -name ilackarms
```

Name not currently implemented but eventually it's going to be the basis of persistence / your "account"
