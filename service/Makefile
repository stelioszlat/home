homecli: main.go internal/server/server.go dashboard/dashboard.go go.mod
	go build -o home
	chmod +x home
	sudo cp home /usr/bin
database: data.go go.mod
	go build -o homedb
	chmod +x homedb
all: homecli

serve:
	go run ./main.go serve

mcp-build:
	go build -o mcp-server ./mcp/server.go
	chmod +x mcp-server
mcp-run:
	./mcp-server

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:zlat123@localhost:5432/home?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:zlat123@localhost:5432/home?sslmode=disable" down

sys: main.c
	gcc -o sysinfo main.c
	chmod +x sysinfo
clean:
	rm -f home homedb sysinfo	