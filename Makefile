homecli: homecli.go homeserver.go go.mod
	go build -o home
	chmod +x home
all: homecli