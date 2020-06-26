FROM openjdk:11

RUN mkdir /app
ADD . /app/
WORKDIR /app

ADD app.jar /app/app.jar
