// Package integration contém os testes de integração de ponta a ponta da API.
//
// Diferente dos testes unitários (que vivem ao lado de cada use case e usam
// mocks isolados), estes testes sobem a aplicação HTTP real — router chi,
// middleware de autenticação, RBAC, handlers e use cases, todos integrados —
// trocando apenas a camada de banco (PostgreSQL) por repositórios falsos em
// memória. Assim cada teste exercita uma rota de verdade, do request HTTP até
// a resposta JSON, sem precisar de banco, Docker ou qualquer passo manual.
//
// Rode com: go test ./internal/integration/...  (ou go test ./... no projeto todo)
package integration
