#!/usr/bin/zsh

git init

go mod init

mkdir cmd
curl -o cmd/main.go https://raw.githubusercontent.com/ezaurum/cthulthu/master/example/cmd/main.go

mkdir docker
mkdir docker/init
curl -o docker/docker-compose.yml https://raw.githubusercontent.com/ezaurum/cthulthu/master/example/docker/docker-compose.yml
curl -o docker/Dockerfile https://raw.githubusercontent.com/ezaurum/cthulthu/master/example/docker/Dockerfile

curl -O https://raw.githubusercontent.com/ezaurum/cthulthu/master/example/modd.conf
curl -O https://raw.githubusercontent.com/ezaurum/cthulthu/master/example/.gitignore
