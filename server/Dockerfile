FROM docker.io/library/golang:1.14.4

RUN git clone https://github.com/ImageMagick/ImageMagick.git
RUN cd ImageMagick/
WORKDIR ImageMagick/
RUN ./configure
RUN make
RUN make install
RUN go get gopkg.in/gographics/imagick.v3/imagick

RUN go get github.com/labstack/echo
RUN go get gopkg.in/ini.v1
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/drive/v3
RUN ldconfig /user/local/lib


RUN mkdir /workdir

WORKDIR  /workdir

COPY . .

RUN apt update