run:
	go build -o chat main.go manager.go user.go groupfunc.go 
	./chat

build:
	go build -o chat main.go manager.go user.go groupfunc.go

clean:
	rm chat