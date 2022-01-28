FROM golang:1.17 AS builder
WORKDIR /go/src/calculator
COPY calculator-app/ .
RUN CGO_ENABLED=0 go build -a calculator_service.go

FROM scratch
COPY --from=builder /go/src/calculator/calculator_service /
CMD ["/calculator_service"]
EXPOSE 8080