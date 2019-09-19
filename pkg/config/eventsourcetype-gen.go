package config

//go:generate go get -u github.com/gertd/gogen-enum
//go:generate gogen-enum -input ./eventsourcetype.yaml -package config -output ./eventsourcetype.go
//go:generate gofmt -w eventsourcetype.go
