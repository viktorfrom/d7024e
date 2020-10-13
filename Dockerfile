# FROM alpine:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go build -o kademlia cmd/main/*.go
RUN go build -o apiclient cmd/client/*.go
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./kademlia"]
