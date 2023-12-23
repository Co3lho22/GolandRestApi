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
