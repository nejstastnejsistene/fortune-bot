fortune-bot
===========

A Slack-bot that wraps around the unix fortune command.

`/fortune --help` prints the man page. Anything else is passed directly to fortune.

Fortunes are stored in `db`, and offensive fortunes are rot-13 encoded in `db/off`. Run `make` to compile newly added fortunes.
