FROM node:24-alpine

ENV NODE_ENV=development
# webpack 4 (vue-cli 4) hashes with md4, which OpenSSL 3 no longer provides
ENV NODE_OPTIONS=--openssl-legacy-provider

ADD start.sh /
RUN chmod +x /start.sh

CMD ["/start.sh"]
