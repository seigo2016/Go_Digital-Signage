FROM arm64v8/debian:buster

RUN apt update && \
    apt install -y --no-install-recommends nodejs npm wget curl

RUN npm config set unsafe-perm true && \
    npm install n -g && \
    n lts

RUN apt purge nodejs npm -y && apt autoremove -y

WORKDIR /workdir
COPY . .

RUN npm install && apt clean -y \
 && rm -rf /var/lib/apt/lists/*
