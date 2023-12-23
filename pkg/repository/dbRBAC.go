package repository

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func GetPermissionsByRoleId(logger *logrus.Logger, db *sql.DB, roleId int) ([]string, error) {
	var permissions []string
	query := "SELECT p.* FROM PERMISSION p INNER JOIN ROLE_PERMISSION rp ON p.id = rp.permission_id WHERE rp.role_id = ?"
	err := db.QueryRow(query, roleId).Scan(&permissions)
	if err != nil {
		logger.WithError(err).WithField("roleId", roleId).Error("Error retrieving permissions using roleId")
		return permissions, err
	}

	logger.WithField("roleId", roleId).Info("Permission retrieved successfully using the roleId ")
	return permissions, nil
}

func GetUserRolesByUserId(logger *logrus.Logger, db *sql.DB, userId int) ([]string, error) {
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
