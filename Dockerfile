FROM golang as builder

COPY . /go

RUN make

FROM registry.suse.com/bci/bci-base
COPY --from=builder /go/soappatrol /bin/soappatrol
