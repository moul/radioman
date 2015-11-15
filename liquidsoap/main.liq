#!/usr/bin/env liquidsoap


# Configure logging
set("log.file", false)
#set("log.file.path", "...")
set("log.stdout", true)
set("log.level", 4)


# Configure Telnet server
#set("server.telnet", true)
#set("server.telnet.port", 23)


# Tuning
set("audio.converter.samplerate.libsamplerate.quality", "best")
set("buffer.path", "/tmp")
set("decoding.buffer_length", 10.)


# Input
input = mksafe(blank())


# Outputs
output.icecast(
  %mp3(bitrate=192),
  host=getenv("ICECAST_PORT_8000_TCP_ADDR"),
  port=int_of_string(getenv("ICECAST_PORT_8000_TCP_PORT")),
  password=getenv("ICECAST_SOURCE_PASSWORD"),
  mount="/radioman.mp3",
  input
)
output.dummy(blank())