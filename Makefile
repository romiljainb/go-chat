run:
	#go build -o chat main.go manager.go user.go groupfunc.go server consts.go
	#go build -o chat main.go manager.go user.go groupfunc.go server/server.go consts.go
	go build -o chat *.go 
	./chat

build:
	go build -o chat main.go manager.go user.go groupfunc.go server consts.go

clean:
	rm chat
