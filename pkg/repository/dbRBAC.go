package repository

import (
	"database/sql"
	"fmt"
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
	//query := "SELECT r.* FROM ROLE r INNER JOIN USER_ROLE ur ON r.id = ur.role_id WHERE ur.user_id = ?"
	query := "SELECT r.name FROM ROLE r INNER JOIN USER_ROLE ur ON r.id = ur.role_id WHERE ur.user_id = ?"
	err := db.QueryRow(query, userId).Scan(&roles)
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error retrieving permissions using roleId")
		return roles, err
	}

	logger.WithField("userId", userId).Info("Roles retrieved successfully using the userId")
	return roles, nil
}

func SetUserRole(logger *logrus.Logger, db *sql.DB, userName string, roleName string) error {

	query := "INSERT INTO USER_ROLE (user_id, role_id) VALUES ((SELECT id FROM USERS WHERE username= ?),(SELECT id FROM ROLE WHERE name= ?))"
	result, err := db.Exec(query, userName, roleName)
	if err != nil {
		logger.WithError(err).WithField("username", userName).Error("Error setting user role for user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).WithField("username", userName).
			Error("Error getting the number of rows affected when trying to store the refreshToken")
		return err
	}

	if !(rowsAffected > 0) {
		logger.WithError(err).WithField("username", userName).
			Warn("No errors setting up the userRole in DB, but rows affected <= 0")
		return fmt.Errorf("no errors setting user role %s for username %s, but no rows affected", roleName, userName)
	}

	logger.WithField("username", userName).Info("Roles for user set up with success")
	return nil

}
