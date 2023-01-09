FROM golang:1.19 as builder

ENV GOOS linux
ENV CGO_ENABLED 0 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app cmd/main.go

FROM alpine:3.14 as production 

RUN apk add --no-cache ca-certificates

COPY --from=builder app . 

EXPOSE 8888

CMD [ "./app" ]

# RUN apt update && apt upgrade -y && \
#     apt install -y git \
#     make openssh-client

# WORKDIR /backend

# RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
#     && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

# CMD air