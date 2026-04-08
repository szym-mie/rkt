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
.PHONY: clean
clean: clean_winnt
else
all: $(TARGET_ELF)
init: $(DR_ELF)
.PHONY: clean
clean: clean_posix
endif

$(DR_EXE) $(DR_ELF):
	@echo -- building $(DR)... --
	cd $(SUB_DR) && $(MAKE)

$(TARGET_EXE): $(DR_EXE) $(PREREQS)
	@echo -- building $@... --
	$(DR_EXE) -res $(RES)
$(TARGET_ELF): $(DR_ELF) $(PREREQS)
	@echo -- building $@... --
	$(DR_ELF) -res $(RES)

.PHONY: clean_winnt
clean_winnt:
	del *.exe
.PHONY: clean_posix
clean_posix:
	rm *.elf