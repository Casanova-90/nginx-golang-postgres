FROM golang:1.13-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

# Instalar dependencias y herramientas adicionales
RUN apk add --no-cache git

# Descargar dependencias del proyecto
RUN go mod tidy

# Comando de depuración para verificar el contenido del directorio
RUN ls -la /go/src/github.com/org/repo

# Construir la aplicación con salida detallada para depuración
RUN go build -v -o server . || (echo "Build failed" && exit 1)

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]