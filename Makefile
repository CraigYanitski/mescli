SHELL = /bin/zsh

.PHONY: help build build_server

help:
	@echo "\n" \
		"                             ___ \n" \
		"   ____ ___  ___  __________/ (_)\n" \
		"  / __ \`__ \/ _ \/ ___/ ___/ / / \n" \
		" / / / / / /  __(__  ) /__/ / /  \n" \
		"/_/ /_/ /_/\___/____/\___/_/_/   \n" \
		"\n" \
		"Use this file to setup MESCLI." \
		"At the moment there is no local storage, so any messages will be deleted after closing the program."

build:
	go build -o mescli ./frontend

build_server:
	go build -o server .
