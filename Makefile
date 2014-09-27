dev:
	grunt clean
	go run server.go

serve:
	grunt updateData
	go build server.go
	./server

get:
	mkdir -p data
	go run get.go $(PAGES)

prod:
	grunt >/dev/null && \
	MARTINI_ENV=production \
	HOST=127.0.0.1 \
	PORT=43434 \
	USERNAME=$(USER) \
	PASSWORD=$(PASSWORD) \
	go run server.go
