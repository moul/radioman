#!/usr/bin/env liquidsoap


# Configure logging
set("log.file", false)
#set("log.file.path", "...")
set("log.stdout", true)
set("log.level", 3)


# Configure Telnet server
set("server.telnet", true)
set("server.telnet.port", 2300)
set("server.telnet.bind_addr", "0.0.0.0")


# Tuning
set("audio.converter.samplerate.libsamplerate.quality", "best")
set("buffer.path", "/tmp")
set("decoding.buffer_length", 10.)


# Telnet variables
radiomand_url = interactive.string("radiomand_url", getenv("RADIOMAND_URL"))


# Helpers
def request_radiomand(path) =
    command = "wget -qO - #{radiomand_url()}/#{path}"
    log("Executing: #{command}")
    get_process_lines(command)
end


# Input
requests = request.queue(id="request")
scheduler = request.dynamic(
    id="scheduler",
    fun () ->
        request.create(
            list.hd(
                request_radiomand("api/liquidsoap/getNextSong")
            )
        )
    )

source = fallback(id="switcher", track_sensitive=true, [
    requests, scheduler, blank(duration=5.)
])


# Outputs
def outputs(s) =
  icy = output.icecast(
    host=getenv("ICECAST_PORT_8000_TCP_ADDR"),
    port=int_of_string(getenv("ICECAST_PORT_8000_TCP_PORT")),
    password=getenv("ICECAST_SOURCE_PASSWORD")
  )

  icy(%mp3(bitrate=128), mount="/mp3-128", s)
  icy(%mp3(bitrate=192), mount="/mp3-192", s)
  icy(%aac(bitrate=128), mount="/aac-128", s)
  icy(%aac(bitrate=192), mount="/aac-192", s)
  icy(%vorbis, mount="/vorbis", s)
  output.dummy(id="manager", s)
end

outputs(source)
