#!/bin/bash

ports=(8080 8081 8082 8083 8084)
for i in ${ports[@]}
do
  ./server_node.sh "$i" &
done

./client.sh ${ports[@]}