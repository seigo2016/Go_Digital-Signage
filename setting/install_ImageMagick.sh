#!/bin/sh

git clone https://github.com/ImageMagick/ImageMagick.git
cd ImageMagick/
./configure
make
make install