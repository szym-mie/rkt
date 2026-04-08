SUB_DR = dev/dev_release
DR = dev_release.out

RES = res
TARGET = rkt
TARGET_EXE = $(TARGET).exe
TARGET_ELF = $(TARGET).elf
PREREQS = main.go $(wildcard src/*.go)

ifeq ($(OS), Windows_NT)
all: $(TARGET_EXE)
.PHONY: clean
clean: clean_winnt
else
all: $(TARGET_ELF)
.PHONY: clean
clean: clean_posix
endif

init: $(DR)

$(DR):
	@echo -- building $@... --
	cd $(SUB_DR) && $(MAKE)

$(TARGET_EXE) $(TARGET_ELF): $(DR) $(PREREQS)
	@echo -- building $@... --
	go get .
	$(DR) -res $(RES)

devel: $(DR)
	@echo -- start devel... --
	go get .
	$(DR) -dev -res $(RES)

.PHONY: clean_winnt
clean_winnt:
	del *.exe 2>nul
	del *.out 2>nul
	del $(RES)\*.zip 2>nul
.PHONY: clean_posix
clean_posix:
	rm *.elf
	rm *.out
	rm $(RES)\*.zip