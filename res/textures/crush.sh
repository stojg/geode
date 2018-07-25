#!/bin/sh -ex


for png in `find . -name "*.png"`;
do
echo "crushing $png"
    convert "$png" -resize 1024x1024 ./temp.png
	pngcrush -rem allb -brute "./temp.png" "$png"
	rm temp.png
done;
