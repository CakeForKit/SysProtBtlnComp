# Стадия сборки (builder) - это "временный слой"
FROM golang:1.24.2-alpine AS builder

RUN mkdir /build
WORKDIR /build

COPY ./go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \                  
    -o /build/bin/SysProtBtlnComp \
    ./cmd/app/main.go


# Финальная стадия - Новый базовый слой
FROM alpine:latest  

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/bin/SysProtBtlnComp .
COPY --from=builder /build/configs ./configs 
COPY --from=builder /build/input_data ./input_data 
COPY --from=builder /build/log ./log 

# Устанавливаем зависимости для runtime (если нужны)
# RUN apk add --no-cache ca-certificates tzdata

CMD ["./SysProtBtlnComp"]

