#!/bin/sh

go install moul.io/generate-fake-data

for voice in Alex Fred Samantha Victoria; do
    text="$(generate-fake-data --dict=phrase --no-stderr --lines=1)"
    say "$text" -v $voice -o $voice-"$(echo $text | tr -cd '[:alnum:]._-')aif"
done
