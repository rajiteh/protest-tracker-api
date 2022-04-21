FROM golang:1.18 as builder

WORKDIR /build

ARG GOOS=linux
ARG GOARCH=amd64

ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0  go build -o ./app main.go

FROM gcr.io/distroless/base-debian11
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]