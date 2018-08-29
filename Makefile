
all: backend front
	go build -o bin/backend backend/main.go 
	go build -o bin/front front/main.go

