# Servidor HTTP/1.1 do zero (Go)

Usa apenas `net` (TCP cru); o HTTP é interpretado em `http1/`.

## Compilar e executar
    go build -o servidor .
    ./servidor --port 8080 --root ./www

- `--port`: porta TCP (padrão 8080), escuta em todas as interfaces
- `--root`: diretório servido (padrão ./www)

## Estrutura
- `main.go`: argumentos e inicialização
- `http1/server.go`: accept, uma goroutine por conexão, keep-alive e timeout ocioso (5s)
- `http1/request.go`: buffer de bytes, parsing da request line e headers, percent-decoding
- `http1/handler.go`: GET/HEAD/405, proteção contra `..` (403), 404
- `http1/response.go`: montagem da resposta, Date (IMF-fixdate), Content-Type
