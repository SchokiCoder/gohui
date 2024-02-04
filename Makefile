VERSION=v0.1.0

hui_d: color.go config.go hui.go menu.go
	go build -o hui_d -gcflags "-l -N" -ldflags "-X 'main.Version=$(VERSION)'"
