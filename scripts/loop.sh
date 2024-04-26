#!/bin/bash

url="http://localhost:17000"

startPos=25
finishPos=100
x=$startPos
y=$startPos
step=10

curl -X POST -d "green" "$url"
curl -X POST  -d "figure $(awk -v s=$startPos 'BEGIN{printf "%.2f %.2f", s/100, s/100}')" "$url"
curl -X POST -d "update" "$url"

sleep 0.01

while true; do
  while ((x < finishPos-startPos)); do
    move_x=$(awk -v s=$step 'BEGIN{printf "%.2f", s/100}')
    curl -X POST -d "move $move_x 0" "$url"
    x=$((x + step))
    curl -X POST -d "update" "$url"
    sleep 0.01
  done

  while ((y < finishPos-startPos)); do
    move_y=$(awk -v s=$step 'BEGIN{printf "%.2f", s/100}')
    curl -X POST -d "move 0 $move_y" "$url"
    y=$((y + step))
    curl -X POST -d "update" "$url"
    sleep 0.01
  done

  while ((x > startPos)); do
    move_x=$(awk -v s=$step 'BEGIN{printf "%.2f", -s/100}')
    curl -X POST -d "move $move_x 0" "$url"
    x=$((x - step))
    curl -X POST -d "update" "$url"
    sleep 0.01
  done

  while ((y > startPos)); do
    move_y=$(awk -v s=$step 'BEGIN{printf "%.2f", -s/100}')
    curl -X POST -d "move 0 $move_y" "$url"
    y=$((y - step))
    curl -X POST -d "update" "$url"
    sleep 0.01
  done
done
