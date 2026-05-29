BEGIN;

CREATE TABLE IF NOT EXISTS log_importacao (
    id                       SERIAL       PRIMARY KEY,
    nome_arquivo             VARCHAR(255) NOT NULL,
    data_inicio              TIMESTAMP    NOT NULL,
    data_fim                 TIMESTAMP,
    usuario_id               VARCHAR(100),
    status                   VARCHAR(50)  NOT NULL,
    total_linhas_processadas INT          NOT NULL DEFAULT 0,
    total_linhas_erro        INT          NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS log_importacao_erro (
    id                SERIAL PRIMARY KEY,
    log_importacao_id INT    NOT NULL REFERENCES log_importacao(id) ON DELETE CASCADE,
    numero_linha      INT    NOT NULL,
    conteudo_original TEXT,
    motivo_erro       TEXT
);

COMMIT;