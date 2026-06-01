BEGIN;

CREATE TABLE IF NOT EXISTS usuario (
    id          SERIAL       PRIMARY KEY,
    email       VARCHAR(255) NOT NULL UNIQUE,
    senha_hash  VARCHAR(255) NOT NULL,
    nome        VARCHAR(100) NOT NULL,
    role        VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'compras')),
    criado_em   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS token_revogado (
    jti        VARCHAR(64)  PRIMARY KEY,
    expira_em  TIMESTAMPTZ  NOT NULL
);

-- Seed: senhas default 'admin123' e 'compras123' (hash bcrypt). Trocar em produção.
INSERT INTO usuario (email, senha_hash, nome, role) VALUES
  ('admin@denarius.local',   '$2a$10$CINdlTtxWY0oBlrnBp15LOZVOW0BX2a8jbO/o2IynJK.VShsoDLbW', 'Administrador', 'admin'),
  ('compras@denarius.local', '$2a$10$UgaRqmMrMTQxq71.onGUnuCDf715PYL0Vg6V0YJgXkCoDJ136m2G2', 'Compras',       'compras')
ON CONFLICT (email) DO NOTHING;

COMMIT;
