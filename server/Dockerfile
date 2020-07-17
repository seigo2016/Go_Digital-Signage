FROM arm64v8/debian:buster

EXPOSE 30000

RUN apt update
RUN apt install git make wget gcc -y && \
    apt clean -y && \
    rm -rf /var/lib/apt/lists/*

RUN wget https://golang.org/dl/go1.14.4.linux-arm64.tar.gz
RUN tar -C /usr/local -xzf go1.14.4.linux-arm64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

RUN go get gopkg.in/ini.v1
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google

RUN go get github.com/gen2brain/go-fitz
RUN go get github.com/nfnt/resize

# RUN go get google.golang.org/api/drive/v3

RUN mkdir /workdir
WORKDIR /workdir
COPY . .

