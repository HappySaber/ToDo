# Используем официальный образ Go
FROM golang:1.23.2 AS builder

RUN mkdir /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы приложения
COPY . .
# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем необходимые библиотеки
RUN apk --no-cache add ca-certificates

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main .
COPY .env .   

# Указываем переменную окружения для порта
ENV PORT=8080

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
