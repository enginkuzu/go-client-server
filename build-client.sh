rm -f go-client
go build -o go-client cmd/client/*
strip go-client
ls -lh go-client
