FROM golang:1.20
WORKDIR /myBotVK
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o myBotVK ./cmd/main.go

ENTRYPOINT ["./myBotVK"]
