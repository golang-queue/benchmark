testing:
	go test -v -run=^$ -benchmem -count=2  -bench .
