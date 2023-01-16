build:
	go test -v -run=^$ -benchmem -count=1  -bench .
