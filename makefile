SUB_DR = dev/dev_release
DR = dev_release
DR_EXE = .\$(DR).exe
DR_ELF = ./$(DR).elf

RES = res
TARGET = rkt
TARGET_EXE = $(TARGET).exe
TARGET_ELF = $(TARGET).elf
PREREQS = main.go $(wildcard src/*.go)

ifeq ($(OS), Windows_NT)
all: $(TARGET_EXE)
init: $(DR_EXE)
dev: dev_winnt
.PHONY: clean
clean: clean_winnt
else
all: $(TARGET_ELF)
init: $(DR_ELF)
dev: dev_posix
.PHONY: clean
clean: clean_posix
endif

$(DR_EXE) $(DR_ELF):
	@echo -- building $(DR)... --
	cd $(SUB_DR) && $(MAKE)

$(TARGET_EXE): $(DR_EXE) $(PREREQS)
	@echo -- building $@... --
	go get .
	$(DR_EXE) -res $(RES)
$(TARGET_ELF): $(DR_ELF) $(PREREQS)
	@echo -- building $@... --
	go get .
	$(DR_ELF) -res $(RES)

dev_winnt: $(DR_EXE)
	@echo -- start devel... --
	go get .
	$(DR_EXE) -dev -res $(RES)
dev_posix: $(DR_ELF)
	@echo -- start devel... --
	go get .
	$(DR_ELF) -dev -res $(RES)

.PHONY: clean_winnt
clean_winnt:
	del *.exe
	del $(RES)\*.zip
.PHONY: clean_posix
clean_posix:
	rm *.elf
	rm $(RES)\*.zip