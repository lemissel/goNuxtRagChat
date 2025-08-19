# Estágio de construção para compilar o binário
FROM golang:1.24-alpine AS builder
WORKDIR /app
ARG TARGETARCH

# Configura a arquitetura de compilação
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${TARGETARCH}


# Copia os arquivos de dependência e baixa
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copia o código-fonte
COPY . .
COPY wait-for-it.sh .
RUN chmod +x ./wait-for-it.sh

# Compila a aplicação e gera o executável 'main'
# Força a compilação para uma arquitetura específica (por exemplo, amd64)
# ENV CGO_ENABLED=0
# ENV GOOS=linux
# ENV GOARCH=amd64
RUN go build -o main -ldflags '-w -s' .

# Estágio de execução, com uma imagem base menor
FROM alpine:latest

# Instala certificados de SSL para garantir que chamadas HTTP/HTTPS funcionem
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copia apenas o executável do estágio de build
# O binário é chamado 'main' e está no diretório /app do estágio builder
COPY --from=builder /app/main .
COPY --from=builder /app/wait-for-it.sh .

# Expõe a porta 8080
EXPOSE 8080

# Define o comando para rodar o executável
CMD ["./main"]