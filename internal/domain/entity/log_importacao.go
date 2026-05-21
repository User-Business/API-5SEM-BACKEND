package entity

import "time"

type LogImportacao struct {
	ID                     int        `json:"id"`
	NomeArquivo            string     `json:"nome_arquivo"`
	DataInicio             time.Time  `json:"data_inicio"`
	DataFim                *time.Time `json:"data_fim,omitempty"`
	UsuarioID              string     `json:"usuario_id"`
	Status                 string     `json:"status"`
	TotalLinhasProcessadas int        `json:"total_linhas_processadas"`
	TotalLinhasErro        int        `json:"total_linhas_erro"`
}

type LogImportacaoErro struct {
	ID              int    `json:"id"`
	LogImportacaoID int    `json:"log_importacao_id"`
	NumeroLinha     int    `json:"numero_linha"`
	ConteudoOriginal string `json:"conteudo_original"`
	MotivoErro      string `json:"motivo_erro"`
}
