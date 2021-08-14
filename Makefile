SHELL=/bin/bash -O extglob -c
generate:
	go run gen/generate_setters/!(*_test).go
