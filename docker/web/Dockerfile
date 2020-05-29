FROM node:13.11.0-alpine3.11

ENV NODE_ENV development

RUN apk add yarn
ADD start.sh /
RUN chmod +x /start.sh

CMD ["/start.sh"]