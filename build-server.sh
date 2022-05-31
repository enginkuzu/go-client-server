rm -f go-server
go build -o go-server cmd/server/*
strip go-server
ls -lh go-server
