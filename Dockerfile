FROM golang:1.22-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache bash git openssh

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файл модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Указываем порт приложения
EXPOSE 8000

# Запускаем main.go
CMD ["go", "run", "cmd/main.go"]
