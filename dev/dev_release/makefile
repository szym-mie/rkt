TARGET = ../../dev_release.out
PREREQS = main.go $(wildcard src/*.go)

all: $(TARGET)

$(TARGET): $(PREREQS)
	go get .
	go build -o $@ .
