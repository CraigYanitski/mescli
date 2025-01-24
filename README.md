# MESCLI
---

MESCLI is a CLI messaging app.
I want secure end-to-end encryption of the messages, as well as a double ratchet for improved security.
There will preferably be markup support, including a maths mode.
It is not yet known whether a server will be made public.
More likely than not, this will have iterations over several stages.

## Development

Since the full project is rather complex, I will focus on a few features for the first stage.
Many of the more complex features making it useful will be deferred to a later release.

- [x] Format text with ANSI codes
- [ ] Generate client- and server-side keys
- [ ] Synchronise client KDFs
- [ ] Encrypt messages before sending through server
- [ ] Decrypt server response
- [ ] Format display -- create TUI
- [ ] ~Format maths env~ DEFERRED
- [ ] ~Allow non-local users~ DEFERRED
- [ ] ~Use DH generation for each message with KDF~ DEFERRED

## Issues

The maths formatting is deferred until I find a reasonable method to format maths in a TUI.
I don't know if I want to make some dependency or attempt to write it myself.

I will also make this functionality as a local experiment for now, deferring any public functionality 
until a much later date.

There are obvious legal reasons for and against full end-to-end encryption on a public service.
I do not want to go through this just for a small project; I just want to test a double-ratchet encryption.
