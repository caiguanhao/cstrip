serve:
	go build server.go
	./server

get:
	mkdir -p data
	go run get.go $(PAGES) > data/commitstrip.json
