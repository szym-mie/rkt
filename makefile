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
else
all: $(TARGET_ELF)
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
