# WebDav Server

## Em Desenvolvimento

Um servidor WebDAV simples com autenticação básica e API para gerenciamento de usuários e diretórios.

### Funcionalidades

- **Criação de Usuários e Diretórios:** Permite criar usuários e seus respectivos diretórios de forma automatizada.
- **Autenticação Básica:** Protege o acesso ao servidor através de nome de usuário e senha.
- **Protocolo WebDAV:** Habilita o acesso e manipulação remota de arquivos através do protocolo WebDAV.
- **API de Gerenciamento:** Oferece endpoints para criar, listar, atualizar e deletar usuários e diretórios.

### Uso

```bash
./dav -root=/var/www/webdav -token=meu_token_secreto
```

**Opções:**

- `-port`: Porta em que o servidor irá escutar (padrão: 8080).
- `-root`: Diretório raiz onde os arquivos serão armazenados (padrão: ./root).
- `-token`: Token de autenticação para a API (padrão: 123456).
- `-i2p`: Habilita o uso do I2P.
- `-config`: Habilita o uso de um arquivo de configuração.
- `-h`: Exibe a ajuda com as opções disponíveis.
- `-export`: Exporta os usuários (./users.json).
- `-import`: Importa os usuários (./users.json).

### API

A API Rest permite gerenciar usuários e diretórios. Para utilizar a API, inclua o token de autenticação no
cabeçalho `Authorization` das requisições:

```
Authorization: Bearer meu_token_secreto
```

**Endpoints:**

- `POST /admin/user`: Cria um novo usuário.
- `GET /admin/user`: Lista todos os usuários.
- `DELETE /admin/user`: Deleta um usuário

Mais informações podem ser obtidas em [Arquivo Http](test.http) ou [Coleção Postman](DavServer.postman_collection.json).

#### Observações

- /user/file e /user/pub O token esperado é um base64 da string `username:password` e o token deve ser passado no header
  da requisição com o nome "Authorization".
- /admin/users O token esperado é um token simples e o token deve ser passado no header da requisição com o nome "
  Authorization" definido no momento da execução do servidor.
- O caminho do DB é por padrão `/tmp/DavServer`, altere o caminho para não perder os dados ao reiniciar o servidor.
- Para Usar com I2P é necessário que o SAM esteja ativo e configurado para o servidor.
- Nenhum usuário é criado por padrão, é necessário criar um usuário para poder acessar o servidor.
- É possível traduzir as mensagens do Srv em `internal/msg`

### Exemplo de Requisição (cURL)

```bash
curl -X POST -H "Authorization: meu_token_secreto" -H "Content-Type: application/json" -d '{"username": "novo_usuario", "password": "senha_forte"}' http://localhost:8080/users
```

### Variáveis de Ambiente
```yaml
# conf.yml

# Configurações principais
APP_NAME: "DavServer"
DB_DIR: "/tmp/DavServer"
PORT: 8080
SHARE_ROOT_DIR: "./root"
TIME_FORMAT: "02-Jan-2006"
TIME_ZONE: "America/Sao_Paulo"
GLOBAL_TOKEN: "123456"
# Configuração Padrão do Servidor
SRV:
  # Caso true removerá a pasta do usuário ao remover o usuário.
  DELETE_FOLDER: false
  CHUNK_SIZE: 500
# Configurações I2P
I2P_CFG:
  ENABLED: false
  HTTP_HOST_AND_PORT: "127.0.0.1:7672"
  URL: "127.0.0.1:7672"
  SAM_ADDR: "127.0.0.1:7656"
  KEY_PATH: "./"
```