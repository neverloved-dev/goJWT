FROM golang:1.23.1

WORKDIR /app

COPY . /app/


RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goJWT

EXPOSE 9000

CMD [ "/goJWT" ]
