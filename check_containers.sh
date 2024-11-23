#!/bin/bash
exclude_image="$1"
servers="$2"
images=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "$servers" | grep -v "$exclude_image")
for image in $images; do
echo "Containers running with image: $image"
id=$(docker ps --filter "ancestor=$image")
docker rm $id -f
docker rmi $image
done
