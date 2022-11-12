.PHONY: run
start: go generate ./ent && go run ../start.go