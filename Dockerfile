FROM golang:1.22.4
LABEL authors="Sheymor"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN mv docs/ cmd/api

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -o /SchoolManagerApi

EXPOSE 8080

# Run
CMD ["/SchoolManagerApi"]
