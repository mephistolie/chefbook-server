package postgres

import "chefbook-server/pkg/logger"

func logRepoError(err error) {
	logger.Error(err.Error())
}
