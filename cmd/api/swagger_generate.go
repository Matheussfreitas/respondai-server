package main

//go:generate sh -c "cd ../.. && go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --outputTypes go,json,yaml"
