# use with GNU make only

SUB_DB = dev/
OUT = rkt
RES = res
PREREQS = main.go $(wildcard src/*.go)

.PHONY: all
all: build
.PHONY: init
init: $(DB)

ifeq ($(OS), Windows_NT)
TARGET = $(OUT).exe
DB = .\dev_bake.out
.PHONY: clean
clean: clean_winnt
else
TARGET = $(OUT).elf
DB = ./dev_bake.out
.PHONY: clean
clean: clean_posix
endif

$(DB):
	@echo -- building $@... --
	cd $(SUB_DB) && $(MAKE)

.PHONY: build
build: $(DB) $(PREREQS)
	@echo -- building $(TARGET)... --
	go get .
	$(DB) -res $(RES) $(FLAGS)

.PHONY: devel
devel: $(DB) $(PREREQS)
	@echo -- start devel $(TARGET)... --
	go get .
	$(DB) -dev -res $(RES) $(FLAGS)

.PHONY: release
release: $(DB) $(PREREQS)
	@echo -- release $(VER)_$(PLAT) --
	@echo -- building $(TARGET)... --
	go get .
	$(DB) -rel -res $(RES) -ver $(VER) -plat $(PLAT) $(FLAGS)

.PHONY: clean_winnt
clean_winnt:
	del *.exe 2>nul
	del *.out 2>nul
	del $(RES)\*.zip 2>nul
.PHONY: clean_posix
clean_posix:
	rm -f *.elf
	rm -f *.out
	rm -f $(RES)\*.zip