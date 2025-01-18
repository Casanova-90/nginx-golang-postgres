# Etapa de construcción
FROM golang:1.13-alpine AS build

# Crear y establecer el directorio de trabajo
WORKDIR /go/src/github.com/Casanova-90/nginx-golang-postgres

# Copiar los archivos de go.mod y go.sum primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto de los archivos del repositorio
COPY . .

# Construir la aplicación Go
RUN go build -o server .

# Etapa final
FROM alpine:3.12

# Exponer el puerto en el que la aplicación escucha
EXPOSE 8000

# Copiar el binario construido desde la etapa de construcción
COPY --from=build /go/src/github.com/Casanova-90/nginx-golang-postgres/server /server

# Comando para ejecutar la aplicación
CMD ["/server"]