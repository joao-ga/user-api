# Use uma imagem oficial do Go como base
FROM golang:1.23.4-alpine

# Defina o diretório de trabalho no contêiner
WORKDIR /app

# Copie os arquivos do projeto para o contêiner
COPY . .

# Baixe as dependências
RUN go mod tidy

# Compile o código Go
RUN go build -o main .

# Exponha a porta da API
EXPOSE 8080

# Comando para rodar a API
CMD ["./main"]
