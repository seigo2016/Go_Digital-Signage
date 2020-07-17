sudo -S apt install wget git unzip -y

mkdir /tmp/ && cd /tmp/

git clone https://github.com/ImageMagick/ImageMagick.git

cd ImageMagick/ && ./configure && make && make install

ldconfig /user/local/lib

cd /home/pi/

wget https://www.dropbox.com/s/bv734uhy5qr85su/release.zip

unzip release.zip

mv release signage-server

cd signage-server && ./main_arm64

# ###################################################### #

apt install wget git unzip -y

git clone https://github.com/ImageMagick/ImageMagick.git

cd ImageMagick/ && ./configure && make && make install

ldconfig /user/local/lib

cd /

wget https://www.dropbox.com/s/bv734uhy5qr85su/release.zip

unzip release.zip

mv release signage-server

cd signage-server && ./main_arm64