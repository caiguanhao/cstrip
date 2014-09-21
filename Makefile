serve:
	go build server.go
	./server

get:
	mkdir -p data
	go run get.go $(PAGES)

prod:
	MARTINI_ENV=production \
	HOST=127.0.0.1 \
	PORT=43434 \
	USERNAME=$(USER) \
	PASSWORD=$(PASSWORD) \
	go run server.go
