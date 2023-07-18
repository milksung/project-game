FROM golang:1.20-alpine

WORKDIR /app

COPY . .
RUN apk update
RUN apk add alpine-sdk
RUN mkdir $HOME/bin 
RUN echo "export PATH=$HOME/bin:$HOME/go/bin:$PATH" >> $HOME/.profile
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.10
RUN swag init
RUN go build -o ./build/API

EXPOSE 3000

CMD [ "./build/API" ]