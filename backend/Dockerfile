FROM golang:1.15-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

# Instalar dependencias
RUN go mod download

# Comando de depuración
RUN ls -la /go/src/github.com/org/repo

# Construir la aplicación
RUN go build -v -o server .

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]
