# WebDav Server

Um servidor WebDAV simples com autenticação básica e API para gerenciamento de usuários e diretórios.

### Funcionalidades

- **Criação de Usuários e Diretórios:** Permite criar usuários e seus respectivos diretórios de forma automatizada.
- **Autenticação Básica:** Protege o acesso ao servidor através de nome de usuário e senha.
- **Protocolo WebDAV:** Habilita o acesso e manipulação remota de arquivos através do protocolo WebDAV.
- **API de Gerenciamento:** Oferece endpoints para criar, listar, atualizar e deletar usuários e diretórios.

### Uso

```bash
./webdav-server -port 8080 -root /var/www/webdav -token meu_token_secreto
```

**Opções:**

- `-port`: Porta em que o servidor irá escutar (padrão: 8080).
- `-root`: Diretório raiz onde os arquivos serão armazenados (padrão: diretório atual).
- `-token`: Token de autenticação para a API (opcional).
- `-h`: Exibe a ajuda com as opções disponíveis.

### API

A API RESTful permite gerenciar usuários e diretórios. Para utilizar a API, inclua o token de autenticação no cabeçalho `Authorization` das requisições:

```
Authorization: Bearer meu_token_secreto
```

**Endpoints:**

- `POST /admin/users`: Cria um novo usuário.
- `GET /admin/users`: Lista todos os usuários.
- `DELETE /users`: Deleta um usuário.


### Exemplo de Requisição (cURL)

```bash
curl -X POST -H "Authorization: meu_token_secreto" -H "Content-Type: application/json" -d '{"username": "novo_usuario", "password": "senha_forte"}' http://localhost:8080/users
```