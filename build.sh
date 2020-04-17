#!/bin/bash
#export GOPHERJS_GOROOT="$(go1.12.16 env GOROOT)"
#GOPHERJS_GOROOT="$(go1.12.16 env GOROOT)" gopherjs build -o static/main.js
GOPHERJS_GOROOT="$(go1.12.16 env GOROOT)" gopherjs build -m -o static/main.min.js