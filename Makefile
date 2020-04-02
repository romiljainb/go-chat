run:
	go build -o chats server/*.go 
	./chats

build:
	go build -o chats server/*.go

clean:
	rm chats
