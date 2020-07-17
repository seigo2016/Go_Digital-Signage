FROM arm64v8/debian:buster

RUN apt update && apt install -y \
    nodejs npm wget curl \
    libgtkextra-dev libgconf2-dev libnss3 libasound2 libxtst-dev libxss1 libgtk-3-0 libgbm-dev
RUN npm config set unsafe-perm true
RUN npm install n -g
RUN n lts

RUN apt purge nodejs npm -y && apt autoremove -y

RUN mkdir /workdir
WORKDIR /workdir
COPY . .

RUN npm install

RUN apt clean -y \
 && rm -rf /var/lib/apt/lists/*

ARG DOCKER_UID=1000
ARG DOCKER_USER=docker
ARG DOCKER_PASSWORD=docker
RUN useradd -m -u ${DOCKER_UID} -G sudo ${DOCKER_USER} \
    && echo ${DOCKER_USER}:${DOCKER_PASSWORD} | chpasswd


