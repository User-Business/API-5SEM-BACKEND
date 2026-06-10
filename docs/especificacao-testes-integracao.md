# Especificação de Testes de Integração — Backend (Nexus / Denarius Data)

> **Autoria:** camada de Testes (Tester).
> **Método:** Dado / Quando / Então (given/when/then), derivado das *user stories* do Product Backlog.
> **Papel deste documento (shift-left):** esta especificação é escrita **antes** do código estar
> finalizado. Ela serve de base para os desenvolvedores praticarem **TDD** (escrever o teste
> unitário a partir do comportamento esperado aqui descrito e só então o código) e para a camada
> de testes implementar os **testes de integração** automatizados que validam a API montada.

> ⚠️ **Importante sobre a stack:** este backend é em **Go (chi + PostgreSQL + JWT)**, com um modelo
> de **Data Warehouse** (tabelas `dim_*` e `fato_*`). Por isso os endpoints e contratos são
> diferentes de um backend CRUD tradicional. Quase **todas** as rotas exigem autenticação (JWT) e
> autorização por papel (RBAC) — esse é o eixo central dos testes.

---

## 1. Rastreabilidade: Casos de Teste → Endpoint → User Story

| Caso de Teste | Endpoint                              | Método | User Story (Backlog)                         |
|---------------|---------------------------------------|--------|----------------------------------------------|
| TC-AUTH01     | `/api/auth/login`                     | POST   | Segurança (RNF — acesso aos dashboards)      |
| TC-AUTH02     | `/api/auth/logout`                    | POST   | Segurança                                    |
| TC-AUTH03     | *(qualquer rota protegida)*           | GET    | Segurança (RBAC)                             |
| TC-PROJ01     | `/api/dim/projetos`                   | GET    | US09 (busca de projetos) / base US02         |
| TC-PRMAT01    | `/api/projetos/materiais`             | GET    | US06 (consumo de material por projeto)       |
| TC-FORN01     | `/api/dim/fornecedores`               | GET    | US09 (busca de fornecedores)                 |
| TC-MAT01      | `/api/dim/materiais`                  | GET    | US09 / US06                                  |
| TC-FCOM01     | `/api/fato/compras`                   | GET    | US02 (custo de projeto)                      |
| TC-FEST01     | `/api/fato/estoque-materiais`         | GET    | US06                                         |
| TC-FEXE01     | `/api/fato/execucao-tarefas`          | GET    | US07 (tempo por tarefa/projeto)              |
| TC-TEMPO01    | `/api/dim/tempo-gasto`                | GET    | US07                                         |
| TC-INV01      | `/api/programa/investimento`          | GET    | US05 (investimento por programa)             |
| TC-PUR01      | `/api/purchases`                      | GET    | US08 (filtros) / US02                        |
| TC-PUR02      | `/api/purchases/metrics`              | GET    | US08 / US02                                  |
| TC-IMP01      | `/api/import-logs`                    | POST/GET | US01 (coleta e organização de dados / ETL) |

**User Stories de referência (Product Backlog):**

- **US01** — Coletar e organizar dados de diferentes sistemas e planilhas em um único lugar (ETL/DW).
- **US02** — Visualizar o custo total de cada projeto.
- **US03** — Visualizar projetos em risco de atraso.
- **US04** — Dashboard correlacionando custo x progresso de execução.
- **US05** — Visualizar qual programa concentra o maior investimento.
- **US06** — Visualizar consumo de material por projeto.
- **US07** — Visualizar tempo gasto por tarefa e projeto.
- **US08** — Filtrar dashboards por programa, projeto, tarefa, material, pedido e período.
- **US09** — Buscar rapidamente projetos, materiais e fornecedores.
- **US10** — Exportar relatórios e dashboards (PDF/CSV) — *sem endpoint no backend ainda*.

---

## 2. Contratos Base

### 2.1 Autenticação (JWT)
Quase todas as rotas ficam sob `/api` e exigem o header:

```
Authorization: Bearer <token-jwt>
```

- **Sem header / token inválido** → `401` com corpo `{"error":"não autorizado"}`.
- **Token válido, mas papel sem permissão** → `403` com corpo `{"error":"acesso negado"}`.
- **Papéis (roles):** `admin` e `compras`.
  - Rotas de dashboard (`/dim/*`, `/fato/*`, `/programa/*`, `/projetos/*`, `/import-logs/*`) → **somente `admin`**.
  - Rotas de compras (`/purchases`, `/purchases/metrics`) → **`admin` ou `compras`**.

### 2.2 Resposta de erro
```json
{ "error": "mensagem" }
```

### 2.3 Lista vazia
> 🔎 **Observação técnica levantada pela camada de testes (shift-left):**
> Hoje, os endpoints `dim_*` e `fato_*` retornam **`null`** quando não há registros (e não `[]`),
> porque serializam um slice nulo. O endpoint `/api/purchases` já trata isso e devolve `[]`.
> **Recomendação para os devs:** padronizar o retorno vazio como `[]` em todos os endpoints de lista.

### 2.4 Login — request/response
Request `POST /api/auth/login`:
```json
{ "email": "string", "senha": "string" }
```
Response `200`:
```json
{ "token": "string", "usuario": { "nome": "string", "role": "admin|compras" } }
```

---

## 3. Cenários — Autenticação e Autorização

### TC-AUTH01 — POST /api/auth/login
**Funcionalidade:** Autenticação de usuário.

```
Cenário: Login com credenciais válidas
  Dado que existe o usuário { email: "admin@empresa.com", senha: "senha123", role: "admin" }
  Quando POST /api/auth/login com {"email":"admin@empresa.com","senha":"senha123"}
  Então o status é 200
  E o body possui "token" não vazio
  E usuario.role == "admin"

Cenário: Senha incorreta
  Quando POST /api/auth/login com {"email":"admin@empresa.com","senha":"errada"}
  Então o status é 401
  E o body é {"error":"credenciais inválidas"}

Cenário: E-mail inexistente
  Quando POST /api/auth/login com {"email":"ninguem@empresa.com","senha":"x"}
  Então o status é 401

Cenário: Corpo malformado (JSON inválido)
  Quando POST /api/auth/login com corpo "{ não é json"
  Então o status é 400
  E o body é {"error":"payload inválido"}
```

### TC-AUTH02 — POST /api/auth/logout
```
Cenário: Logout revoga o token (denylist)
  Dado um usuário autenticado com token válido T
  Quando POST /api/auth/logout com Authorization: Bearer T
  Então o status é 200
  E ao reutilizar T em qualquer rota protegida, o status passa a ser 401
```

### TC-AUTH03 — Controle de acesso (RBAC) em rotas protegidas
```
Cenário: Acesso sem token
  Quando GET /api/dim/projetos sem header Authorization
  Então o status é 401

Cenário: Acesso autenticado, mas sem permissão
  Dado um usuário com role "compras"
  Quando GET /api/dim/projetos (rota exclusiva de admin)
  Então o status é 403

Cenário: Acesso autorizado
  Dado um usuário com role "admin"
  Quando GET /api/dim/projetos
  Então o status é 200
```

---

## 4. Cenários — Projetos e Dimensões

### TC-PROJ01 — GET /api/dim/projetos  (US09 / base US02)
**Funcionalidade:** Listagem da dimensão de projetos.

```
Cenário: Listar projetos (admin)
  Dado que existem projetos na dimensão dim_projeto
  Quando GET /api/dim/projetos com token de admin
  Então o status é 200
  E Content-Type é application/json
  E o body é uma lista onde cada item possui:
      sk_projeto, id_projeto, codigo_projeto, nome_projeto, responsavel,
      status, codigo_programa, nome_programa, gerente_programa,
      data_inicio, data_fim_prevista   (todos string)

Cenário: Falha na camada de dados
  Dado que o repositório retorna erro
  Quando GET /api/dim/projetos com token de admin
  Então o status é 500
```

### TC-FORN01 — GET /api/dim/fornecedores  (US09)
```
Cenário: Listar fornecedores (admin)
  Quando GET /api/dim/fornecedores com token de admin
  Então o status é 200
  E cada item possui: sk_fornecedor, id_fornecedor, codigo_fornecedor,
      razao_social, cidade, estado, categoria, status
```

### TC-MAT01 — GET /api/dim/materiais  (US09 / US06)
```
Cenário: Listar materiais (admin)
  Quando GET /api/dim/materiais com token de admin
  Então o status é 200
  E cada item possui: sk_material, id_material, codigo_material,
      descricao, categoria, fabricante, custo_estimado, status
```

### TC-PRMAT01 — GET /api/projetos/materiais  (US06)
**Funcionalidade:** Consumo de material por projeto.

```
Cenário: Retornar materiais por projeto (admin)
  Dado que há vínculos projeto↔material no banco
  Quando GET /api/projetos/materiais com token de admin
  Então o status é 200
  E cada item possui: codigo_projeto (string), nome_projeto (string),
      codigo_material (string), descricao_material (string),
      quantidade_estoque (integer)
```

---

## 5. Cenários — Tabelas Fato (análises)

### TC-FCOM01 — GET /api/fato/compras  (US02 — custo)
```
Cenário: Retornar fatos de compras (admin)
  Quando GET /api/fato/compras com token de admin
  Então o status é 200
  E cada item possui: sk_fato, sk_projeto, sk_fornecedor, sk_solicitacao,
      sk_tempo, valor_total_pedido, valor_alocado_projeto   (todos string)
```

### TC-FEST01 — GET /api/fato/estoque-materiais  (US06)
```
Cenário: Retornar fatos de estoque (admin)
  Quando GET /api/fato/estoque-materiais com token de admin
  Então o status é 200
  E cada item possui: sk_fato, sk_projeto, sk_material, sk_tempo,
      quantidade_estoque, quantidade_empenhada
```

### TC-FEXE01 — GET /api/fato/execucao-tarefas  (US07)
```
Cenário: Retornar fatos de execução de tarefas (admin)
  Quando GET /api/fato/execucao-tarefas com token de admin
  Então o status é 200
  E cada item possui: sk_fato, sk_projeto, sk_tarefa, sk_responsavel,
      sk_tempo, horas_trabalhadas
```

### TC-TEMPO01 — GET /api/dim/tempo-gasto  (US07)
```
Cenário: Retornar tempo gasto consolidado (admin)
  Quando GET /api/dim/tempo-gasto com token de admin
  Então o status é 200
  E o body é JSON (estrutura agregada de horas por tarefa/projeto)
```

---

## 6. Cenários — Programa (Investimento) — US05

### TC-INV01 — GET /api/programa/investimento
**Funcionalidade:** Investimento total agrupado por programa.

```
Cenário: Retornar investimento por programa (admin)
  Dado que existe o programa { codigo_programa: "PRG-01",
                               nome_programa: "Programa Alfa",
                               investimento_total: 1500.50 }
  Quando GET /api/programa/investimento com token de admin
  Então o status é 200
  E o body contém um item com:
      codigo_programa == "PRG-01", nome_programa == "Programa Alfa",
      investimento_total == 1500.50
  E cada item possui: codigo_programa (string), nome_programa (string),
      investimento_total (number)

Cenário: Acesso negado para role "compras"
  Quando GET /api/programa/investimento com token de "compras"
  Então o status é 403
```

---

## 7. Cenários — Compras (Purchases) — US08 / US02

### TC-PUR01 — GET /api/purchases
**Funcionalidade:** Listagem de solicitações (SC) e pedidos (PC) de compra, com filtros.
Aceita os filtros via query string: `type` (SC|PC), `status`, `start_date`, `end_date`.

```
Cenário: Listar compras (admin ou compras)
  Quando GET /api/purchases com token de admin
  Então o status é 200
  E o body é uma lista (ou [] quando vazio)
  E cada item possui: id (integer), type ("SC"|"PC"), numero (string),
      status (string), data_criacao (string),
      data_previsao_entrega (string, opcional), duracao_dias (integer, opcional),
      atrasado (boolean, opcional)

Cenário: Perfil "compras" também tem acesso
  Quando GET /api/purchases com token de role "compras"
  Então o status é 200

Cenário: Filtrar por tipo
  Quando GET /api/purchases?type=PC com token de admin
  Então o status é 200
  E todos os itens possuem type == "PC"

Cenário: Acesso sem token
  Quando GET /api/purchases sem Authorization
  Então o status é 401
```

### TC-PUR02 — GET /api/purchases/metrics
```
Cenário: Retornar métricas agregadas de compras
  Quando GET /api/purchases/metrics com token de admin
  Então o status é 200
  E o body possui: total_purchases (integer), total_sc (integer),
      total_pc (integer), average_duration_pc (number),
      status_counts (lista de {status, count})
```

---

## 8. Cenários — Logs de Importação (ETL) — US01

### TC-IMP01 — /api/import-logs
**Funcionalidade:** Registro do resultado das importações de dados (ETL).

```
Cenário: Criar log de importação (admin)
  Quando POST /api/import-logs com um corpo de log válido e token de admin
  Então o status é 200/201
  E o body retorna o id criado

Cenário: Listar logs de importação (admin)
  Quando GET /api/import-logs com token de admin
  Então o status é 200
  E o body é uma lista de logs

Cenário: Acesso sem permissão
  Quando GET /api/import-logs com token de role "compras"
  Então o status é 403
```

---

## 9. Dados de Seed Necessários

Para reproduzir os cenários acima de forma determinística, o ambiente de teste precisa de:

| Entidade                | Quantidade mínima | Observação                                                |
|-------------------------|-------------------|-----------------------------------------------------------|
| Usuários                | 2                 | Um `admin` e um `compras`, com senha conhecida (bcrypt)   |
| Programas               | 2                 | Com investimentos diferentes (para US05)                  |
| Projetos (dim)          | 2+                | Vinculados a programas distintos                          |
| Materiais (dim)         | 2+                | Para US06/US09                                            |
| Fornecedores (dim)      | 2+                | Para US09                                                 |
| Fatos de compras        | 2+                | Valores distintos (US02)                                  |
| Fatos de estoque        | 2+                | Para US06                                                 |
| Fatos de execução       | 2+                | Horas trabalhadas distintas (US07)                        |
| Compras (SC/PC)         | 2+                | Tipos e status variados, para filtros (US08)              |

> **Nos testes de integração automatizados (Go + httptest):** o "banco" é substituído por
> repositórios falsos em memória, que recebem exatamente esses dados de seed por código.
> Assim cada teste prepara e limpa seus próprios dados, sem depender de um PostgreSQL real.

---

## 10. Como os testes automatizados implementam esta especificação

Os testes de integração ficam em `internal/integration/` e sobem a aplicação HTTP real
(router chi + middleware de auth + RBAC + handlers + use cases), trocando apenas a camada de
banco por repositórios em memória. Rodam com:

```bash
go test ./...                          # tudo (unitário + integração)
go test -v ./internal/integration/...  # apenas integração, detalhado
```

Cada subteste é nomeado com o código do caso (ex.: `TC-AUTH01`), garantindo rastreabilidade
direta entre esta especificação e o código de teste.
