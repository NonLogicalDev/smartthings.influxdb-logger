#######################################################################
#                                BUILD                                #
#######################################################################

FROM golang:1.14 AS builder

WORKDIR /build/src/smartthings.logger
COPY . .

RUN CGO_ENABLED=0 go build -o /build/smt.logger ./cmd/smt.logger

#######################################################################
#                                 RUN                                 #
#######################################################################

FROM alpine
COPY --from=builder /build/smt.logger /usr/bin/smt.logger

ENV SMT_LISTEN=0.0.0.0:5555
ENV SMT_INFLUX_URL=http://0.0.0.0:8086

CMD ["smt.logger"]
