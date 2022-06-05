FROM golang:latest

COPY . /workspace
WORKDIR /workspace
CMD "/bin/sh"