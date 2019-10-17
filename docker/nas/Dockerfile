FROM golang:buster
#FROM debian:stretch

#RUN bash -c 'cat /etc/apt/sources.list'

#RUN printf "deb http://httpredir.debian.org/debian jessie-backports main non-free\ndeb-src http://httpredir.debian.org/debian jessie-backports main non-free" > /etc/apt/sources.list.d/backports.list

RUN apt-get update && apt-get install --no-install-recommends build-essential gcc-arm-linux-gnueabihf -y
