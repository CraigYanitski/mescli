# MESCLI
---

## Introduction

MESCLI will be a CLI messaging app, and is currently in development.
I want secure end-to-end encryption of the messages, as well as a 
double-ratchet for improved security.
There will preferably be markup support, eventually including a maths mode.
It is not yet known whether a server will be made public (obviously a public 
server is necessary to make this functional beyond testing).
More likely than not, this will have many iterations over several stages.

## Motivation

I am fascinated by cryptology, so I wanted to attempt my own implementation of 
an end-to-end messaging service.
For the double-ratchet algorithm I am implementing for this app, I am using as 
a guide the very thorough [Signal protocol](https://signal.org/docs/specifications/doubleratchet/) 
specifications.
One can also find there a [description](https://signal.org/docs/specifications/x3dh/) 
of the extended triple-Diffie-Hellman (X3DH) asynchronous key exchange that is 
also implemented in `mescli`.

## Installing / Building

I currently have two executables you can build: one for the server and the 
other for the client.
Since MESCLI is actively developed and not yet deployed, this defaults to 
serving on the local host.
You can compile and run the client executable (called `mescli`) using,

```bash
make build
./mescli
```

It uses a configuration file from the current directory (this will change soon).
This file will hold the user's identifying keys as well as the ratchet keys and 
their conversations.
In order to be cryptographically secure, messages are not stored on the server.
There is not much to test now other than creating an account on the server and 
initialising your keys.
I am currently writing handlers for the `messages` endpoint.

The `server` executable can also be compiled and run using,

```bash
make build-server
./server
```

The server interacts with a PostgreSQL database, and the required migrations
are stored in `./sql/schema` and run using [`goose`](https://github.com/pressly/goose).

In order for the server to run correctly, it expects a `.env` file in the root directory
with environment variables `DB_URL` (defined by the PostgreSQL database) and 
`JWT_SECRET`, a base 64 encoded string used for authenticating users.
You can generate your own 32-byte key using,

```bash
openssl rand -base64 32
```

## Development

Since the full project is rather complex, I will focus on a few features for 
the first stage.
Many of the more complex features making it useful will be deferred to a later 
release.

- [x] Format text with ANSI codes
  - [x] replace with glamour/lipgloss
- [x] Generate client- and server-side keys
- [x] Synchronise client KDFs
- [x] Encrypt messages before sending through server
- [x] Decrypt server response
- [x] Format display
- [ ] ~Use DH generation for each message with KDF~ DEFERRED until after server setup
- [x] Create a TUI
  - [x] create TUI framework (bubbletea)
  - [x] implement Markdown rendering
  - [ ] implement encrypted messaging
- [ ] ~Allow non-local users~ DEFERRED (might deploy on AWS)
  - [x] develop server code
- [ ] ~Format maths env~ DEFERRED to different project extending goldmark
  - [ ] ~use MathJax with glamour~ DEFERRED

## Issues

The maths formatting is deferred until I find a reasonable method to format maths in a TUI.
I don't know if I want to make some dependency or attempt to write it myself.

I will also make this functionality as a local experiment for now, deferring any public functionality 
until a much later date.

There are obvious legal reasons for and against full end-to-end encryption on a public service.
I do not want to go through this just for a small project; I just want to test a double-ratchet encryption.

## Contributing

Consider contacting me before trying to contribute.
I am happy to review any improvements you devise.

### clone the repo

```bash
git clone https://github.com/CraigYanitski/mescli@latest
cd mescli
```

Then implement a new process, create some tests for your contribution, and submit a pull request.
I should be able to respond within a day :-)

