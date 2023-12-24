package repository

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

// GetPermissionsByRoleId retrieves a list of permissions associated with a specific role ID from the database.
//
// logger: A logrus.Logger instance for logging information and errors.
// db: A pointer to the SQL database connection.
// roleId: The ID of the role for which permissions are to be retrieved.
//
// Returns a slice of strings containing permission names and an error, if any.
// If successful, the permissions are retrieved and returned without errors.
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

// GetUserRolesByUserId retrieves a list of role names associated with a specific user ID from the database.
//
// logger: A logrus.Logger instance for logging information and errors.
// db: A pointer to the SQL database connection.
// userId: The ID of the user for which roles are to be retrieved.
//
// Returns a slice of strings containing role names and an error, if any.
// If successful, the roles are retrieved and returned without errors.
func GetUserRolesByUserId(logger *logrus.Logger, db *sql.DB, userId int) ([]string, error) {
	var roles []string
	query := "SELECT r.name FROM ROLE r INNER JOIN USER_ROLE ur ON r.id = ur.role_id WHERE ur.user_id = ?"

	rows, err := db.Query(query, userId)
	if err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error retrieving roles using userId")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			logger.WithError(err).WithField("userId", userId).Error("Error scanning role")
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		logger.WithError(err).WithField("userId", userId).Error("Error iterating over roles")
		return nil, err
	}

	logger.WithField("userId", userId).Info("Roles retrieved successfully using the userId")
	return roles, nil
}

// SetUserRole associates a user with a specific role in the database.
//
// logger: A logrus.Logger instance for logging information and errors.
// db: A pointer to the SQL database connection.
// userName: The username of the user to whom the role should be assigned.
// roleName: The name of the role to be assigned to the user.
//
// Returns an error if there is any issue while setting the user role.
// If successful, the user is associated with the specified role without errors.
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
