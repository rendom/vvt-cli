.PHONY: all clean build

DIST := ./build

all: clear build
build: 
	@GO15VENDOREXPERIMENT=1 go build -o $(DIST)/vvt

clear:
	$(RM) -r $(DIST)/* 
