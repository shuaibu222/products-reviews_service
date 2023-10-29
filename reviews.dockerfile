FROM alpine:latest

RUN mkdir /app

COPY reviewsApp /app

CMD [ "/app/reviewsApp"]