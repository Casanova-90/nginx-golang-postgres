FROM golang:1.15-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

# Instalar dependencias y herramientas adicionales
RUN apk add --no-cache git

# Descargar dependencias del proyecto
RUN  go mod tidy -v

# Comando de depuración para verificar el contenido del directorio
RUN ls -la /go/src/github.com/org/repo

# Construir la aplicación con salida detallada para depuración y redirigir la salida de error a un archivo
RUN go build -v -o server . > build.log 2>&1 || (cat build.log && exit 1)

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]