FROM golang:1.16-alpine AS builder

WORKDIR /work

ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o yaml-merge main.go
RUN chmod +x yaml-merge

FROM alpine 
COPY --from=builder /work/yaml-merge /usr/bin
# Verify binary can run
RUN yaml-merge -h
ENTRYPOINT [ "yaml-merge" ]
