package postgres

import "github.com/mephistolie/chefbook-server/pkg/logger"

func logRepoError(err error) {
	logger.Error(err.Error())
}
