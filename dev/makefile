TARGET = ../dev_bake.out
PREREQS = main.go $(wildcard src/*.go)

all: $(TARGET)

$(TARGET): $(PREREQS)
	go get .
	go build -o $@ .
