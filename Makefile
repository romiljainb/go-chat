run:
	go build -o chat main.go manager.go user.go groupfunc.go server.go consts.go
	./chat

build:
	go build -o chat main.go manager.go user.go groupfunc.go server.go consts.go

clean:
	rm chat