build:
	go test -v -run=^$ -benchmem -count=5  -bench .
