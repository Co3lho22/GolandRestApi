package repository

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func GetRolesByUserId(logger *logrus.Logger, db *sql.DB, userId int) ([]string, error) {
	var roles []string
	query := "SELECT r.* FROM ROLE r INNER JOIN USER_ROLE ur ON r.id = ur.role_id WHERE ur.user_id = ?"
	err := db.QueryRow(query, userId).Scan(&roles)
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error retrieving permissions using roleId")
		return roles, err
	}

	logger.WithField("userId", userId).Info("Roles retrieved successfully using the userId")
	return roles, nil
}
