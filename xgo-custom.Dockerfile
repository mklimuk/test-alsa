from karalabe/xgo-latest

RUN dpkg --add-architecture armhf && cd /tmp && \
    echo "deb [arch=armhf] http://ports.ubuntu.com/ xenial main universe" > /etc/apt/sources.list && \
    apt-get update && apt-get download libasound2:armhf libasound2-dev:armhf && \
    dpkg -i --ignore-depends=libc6:armhf libasound2_*.deb libasound2-dev_*.deb && \
    rm libasound2*.deb
