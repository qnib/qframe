FROM qnib/uplain-golang

ARG LIBRDKAFKA_VER=0.11.0.x
RUN apt-get -qq update \
 && apt-get -qq install -y wget bsdtar python librdkafka-dev \
 && cd /usr/local/src/ \
 && wget -qO -  https://github.com/edenhill/librdkafka/archive/${LIBRDKAFKA_VER}.zip| bsdtar xzf - -C . \
 && cd librdkafka-${LIBRDKAFKA_VER} \
 && chmod +x configure lds-gen.py \
 && ./configure \
 && make \
 && make install

