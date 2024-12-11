# syntax=docker/dockerfile:1

# Etapa de build
FROM golang:1.22 AS builder

# Defina o diretório de trabalho no contêiner
WORKDIR /app

# Copie arquivos de dependências
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod download

# Copie o restante dos arquivos
COPY . .

# Compile a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

# Etapa de execução
FROM gcr.io/distroless/static:nonroot

# Defina o diretório de trabalho
WORKDIR /app

# Copie o binário da etapa de build
COPY --from=builder /app/main .

# Exponha a porta usada pela aplicação
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./main"]
