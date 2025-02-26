SHELL = /bin/zsh

.PHONY: help build build_server

help:
	@cat assets/logo.txt
	@echo "Use this file to setup MESCLI." \
		"At the moment there is no local storage, so any messages will be deleted after closing the program."

build:
	go build -o mescli ./frontend
	@echo "\n\n\n"
	@cat assets/logo.txt

build_server:
	go build -o server .
	@echo "\n\n\n"
	@cat assets/logo.txt
