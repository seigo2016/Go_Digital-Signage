FROM arm64v8/debian:buster

EXPOSE 30000

RUN apt update && apt install wget git gcc ca-certificates -y && \
    apt clean -y && \
    rm -rf /var/lib/apt/lists/*

RUN wget  --no-check-certificate https://golang.org/dl/go1.14.4.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.14.4.linux-arm64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

RUN go get gopkg.in/ini.v1 \
    golang.org/x/oauth2 \
    golang.org/x/oauth2/google \
    github.com/gen2brain/go-fitz \
    github.com/nfnt/resize

WORKDIR /workdir
COPY . .

