#!/bin/zsh
docker build -t cav-test .

docker run --rm -it --name=cav-01 \
  -v $HOME/dev/go/cav/clamav-db:/var/lib/clamav:rw \
  -v $HOME/dev/go/cav/clamav-scan:/temp/scan:rw \
  -p 8080:8080 \
  cav-test