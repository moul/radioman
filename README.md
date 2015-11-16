# :musical_note: radioman

[![GoDoc](https://godoc.org/github.com/moul/radioman/radioman?status.svg)](https://godoc.org/github.com/moul/radioman/radioman)

![Logo](https://raw.githubusercontent.com/moul/radioman/master/radioman/web/static/radioman.png)

Web radio solution using Liquidsoap and Icecast

## Screenshots

:warning: **WORK IN PROGRESS** :warning:

![](https://raw.githubusercontent.com/moul/radioman/master/assets/screenshot-001.png)

## Run using Docker

Requires `docker` and `docker-compose`.

```bash
# Clone the repository
git clone github.com/moul/radioman
cd radioman

# Run Docker containers using docker-compose
make compose

# Open in the browser (adjust 'localhost' to your Docker host ip address if needed)
$BROWSER http://localhost:4343/
$BROWSER http://localhost:4343/admin
$PLAYER http://localhost:4343/playlist.m3u
```

## License
MIT
