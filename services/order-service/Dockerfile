# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy full project (IMPORTANT for proto + internal imports)
COPY . .

# download dependencies
RUN go mod download

# build binary
RUN go build -o order-service ./services/order-service/cmd/main.go


# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/order-service .

EXPOSE 50051

CMD ["./order-service"]