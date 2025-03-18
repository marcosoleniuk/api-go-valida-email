# API de Validação de E-mails

## Descrição
Esta é uma API projetada para realizar a validação de e-mails, fornecendo uma resposta detalhada sobre a validade e outras informações relacionadas.

## Versão
- **1.0**

[//]: # (## Link de da API em produção)

[//]: # (- **[https://api-email-validation.herokuapp.com/api]&#40;https://api-email-validation.herokuapp.com/api&#41;**)

## Base Path
```
/api
```

## Endpoints

### 1. Valida Múltiplos E-mails
**Rota:** `POST /emails`

**Descrição:** Valida uma lista de e-mails fornecida no corpo da requisição.

#### Parâmetros
- **Body:**
```json
{
  "emails": ["exemplo1@email.com", "exemplo2@email.com"]
}
```
- **Header:**
```
X-API-KEY: {sua_chave_de_autenticacao}
```

#### Respostas
- **200 OK**
```json
[
  {
    "domain": "email.com",
    "email": "exemplo1@email.com",
    "error": "",
    "is_valid": true,
    "validation_type": "syntax"
  }
]
```
- **400 Bad Request**
```json
{
  "error": "Corpo da requisição inválido"
}
```

---

### 2. Valida um Único E-mail
**Rota:** `GET /valida-email`

**Descrição:** Valida um e-mail fornecido via query string.

#### Parâmetros
- **Query:** `email={email_para_validacao}`
- **Header:**
```
X-API-KEY: {sua_chave_de_autenticacao}
```

#### Respostas
- **200 OK**
```json
{
  "domain": "email.com",
  "email": "exemplo@email.com",
  "error": "",
  "is_valid": true,
  "validation_type": "syntax"
}
```
- **400 Bad Request**
```json
{
  "error": "E-mail inválido"
}
```
- **500 Internal Server Error**
```json
{
  "error": "Erro interno no servidor"
}
```

---

## Definições de Esquema

### **main.BulkEmailRequest**
- **`emails`** (array de strings) - Lista de e-mails para validação.

### **main.ValidationResponse**
- **`domain`** (string) - Domínio do e-mail validado.
- **`email`** (string) - O e-mail que foi validado.
- **`error`** (string) - Mensagem de erro caso a validação falhe.
- **`is_valid`** (boolean) - Indica se o e-mail é válido.
- **`validation_type`** (string) - Tipo de validação realizada (ex.: sintaxe, MX, etc.).

---

## Autenticação
Esta API utiliza autenticação via chave de API. Inclua a seguinte chave no cabeçalho das requisições:
```
X-API-KEY: {lkWVR82OF5Ffo3VjVlL8GRkXPn5dt261gq7wOfxDUIokVE0nnGcJ5EN1NeNbUliQ}
```
Ou no Parametro de URL:
```
/valida-email?apiKey=lkWVR82OF5Ffo3VjVlL8GRkXPn5dt261gq7wOfxDUIokVE0nnGcJ5EN1NeNbUliQ
```

