# Build prover-server
FROM golang:1.19-bullseye as base

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN go build -o ./prover ./cmd/prover/prover.go
RUN go build -tags="rapidsnark_noasm" -o ./prover_noasm ./cmd/prover/prover.go


# Main image
FROM ubuntu:22.04
RUN apt update
RUN apt install -y g++
RUN apt-get install -y nlohmann-json3-dev
RUN apt install -y libmpc-dev
RUN apt-get -y install nasm

COPY --from=base /build/prover /home/app/prover
COPY --from=base /build/prover_noasm /home/app/prover_noasm
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY docker-entrypoint.sh /usr/local/bin/

COPY ./configs   /home/app/configs
COPY ./circuits  /home/app/circuits

WORKDIR /home/app

# Command to run
ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 8000
