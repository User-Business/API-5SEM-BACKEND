CREATE TABLE IF NOT EXISTS programa (
    id_programa        SERIAL       PRIMARY KEY,
    codigo_programa    VARCHAR(50)  NOT NULL UNIQUE,
    nome_programa      VARCHAR(100) NOT NULL,
    gerente_programa   VARCHAR(100),
    gerente_tecnico    VARCHAR(100),
    data_inicio        DATE,
    data_fim_prevista  DATE,
    status             VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS projeto (
    id_projeto         SERIAL       PRIMARY KEY,
    codigo_projeto     VARCHAR(50)  NOT NULL UNIQUE,
    nome_projeto       VARCHAR(100) NOT NULL,
    id_programa        INT          REFERENCES programa(id_programa),
    responsavel        VARCHAR(100),
    custo_hora         DECIMAL(10,2),
    data_inicio        DATE,
    data_fim_prevista  DATE,
    status             VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS tarefa (
    id_tarefa          SERIAL       PRIMARY KEY,
    codigo_tarefa      VARCHAR(50)  NOT NULL UNIQUE,
    id_projeto         INT          NOT NULL REFERENCES projeto(id_projeto),
    titulo             VARCHAR(200) NOT NULL,
    responsavel        VARCHAR(100),
    estimativa_horas   INT,
    data_inicio        DATE,
    data_fim_prevista  DATE,
    status             VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS tempo_tarefa (
    id_tempo           SERIAL       PRIMARY KEY,
    id_tarefa          INT          NOT NULL REFERENCES tarefa(id_tarefa),
    usuario            VARCHAR(100),
    data               DATE         NOT NULL,
    horas_trabalhadas  DECIMAL(5,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS material (
    id_material        SERIAL       PRIMARY KEY,
    codigo_material    VARCHAR(50)  NOT NULL UNIQUE,
    descricao          VARCHAR(200) NOT NULL,
    categoria          VARCHAR(100),
    fabricante         VARCHAR(100),
    custo_estimado     DECIMAL(10,2),
    status             VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS fornecedor (
    id_fornecedor      SERIAL       PRIMARY KEY,
    codigo_fornecedor  VARCHAR(50)  NOT NULL UNIQUE,
    razao_social       VARCHAR(200) NOT NULL,
    cidade             VARCHAR(100),
    estado             CHAR(2),
    categoria          VARCHAR(100),
    status             VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS solicitacao_compra (
    id_solicitacao       SERIAL       PRIMARY KEY,
    numero_solicitacao   VARCHAR(50)  NOT NULL UNIQUE,
    id_projeto           INT          NOT NULL REFERENCES projeto(id_projeto),
    id_material          INT          NOT NULL REFERENCES material(id_material),
    quantidade           INT          NOT NULL,
    data_solicitacao     DATE         NOT NULL,
    prioridade           VARCHAR(50),
    status               VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS pedido_compra (
    id_pedido            SERIAL       PRIMARY KEY,
    numero_pedido        VARCHAR(50)  NOT NULL UNIQUE,
    id_solicitacao       INT          NOT NULL REFERENCES solicitacao_compra(id_solicitacao),
    id_fornecedor        INT          NOT NULL REFERENCES fornecedor(id_fornecedor),
    data_pedido          DATE         NOT NULL,
    data_previsao_entrega DATE,
    valor_total          DECIMAL(12,2),
    status               VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS compra_projeto (
    id_compra_projeto    SERIAL       PRIMARY KEY,
    id_pedido            INT          NOT NULL REFERENCES pedido_compra(id_pedido),
    id_projeto           INT          NOT NULL REFERENCES projeto(id_projeto),
    valor_alocado        DECIMAL(12,2)
);

CREATE TABLE IF NOT EXISTS empenho (
    id_empenho           SERIAL       PRIMARY KEY,
    id_projeto           INT          NOT NULL REFERENCES projeto(id_projeto),
    id_material          INT          NOT NULL REFERENCES material(id_material),
    quantidade_empenhada INT          NOT NULL,
    data_empenho         DATE         NOT NULL
);

CREATE TABLE IF NOT EXISTS estoque_projeto (
    id_estoque           SERIAL       PRIMARY KEY,
    id_projeto           INT          NOT NULL REFERENCES projeto(id_projeto),
    id_material          INT          NOT NULL REFERENCES material(id_material),
    quantidade           INT          NOT NULL DEFAULT 0,
    localizacao          VARCHAR(100)
);