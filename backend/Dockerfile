FROM golang:1.15-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

# Instalar dependencias y herramientas adicionales
RUN apk add --no-cache git
RUN go mod tidy

# Comando de depuración
RUN ls -la /go/src/github.com/org/repo

# Construir la aplicación con salida de depuración
RUN go build -v -o server . || (echo "Build failed" && exit 1)

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]