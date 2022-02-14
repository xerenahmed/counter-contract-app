all:
	go build -o build/current current.go
	go build -o build/increment increment.go

