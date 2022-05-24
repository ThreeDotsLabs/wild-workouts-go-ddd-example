FROM node:12.22.8-alpine3.13

ENV NODE_ENV development

RUN apk --no-cache add yarn python2 make gcc g++
ADD start.sh /
RUN chmod +x /start.sh

CMD ["/start.sh"]