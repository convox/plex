# convox/plex

TCP multiplexing over a shell-based stream

## Installation

Add the static server to your container image:

    ADD https://convox-binaries.s3.amazonaws.com/plex-static /plex-static

Install the client on your local machine

    $ go get -u github.com/convox/plex

## Usage

Use the client to multiplex a TCP stream over a one-off shell:

    $ plex -l 4369:4369 -r 4369:4369 convox exec PID /plex-static server

Specify `-l 3000:4000` to forward local port `3000` to port `4000` inside the container.

Specify `-r 5000:6000` to forward port `5000` inside the container to port `6000` locally.

You can specify both of these options as many times as you like.
