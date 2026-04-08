TARGET = ../../dev_release
TARGET_EXE = $(TARGET).exe
TARGET_ELF = $(TARGET).elf
PREREQS = main.go $(wildcard src/*.go)

ifeq ($(OS), Windows_NT)
all: $(TARGET_EXE)
else
all: $(TARGET_ELF)
endif

$(TARGET_EXE): $(PREREQS)
	go get .
	go build -o $@ .
$(TARGET_ELF): $(PREREQS)
	go get .
	go build -o $@ .
