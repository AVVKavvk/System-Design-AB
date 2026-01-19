package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AVVKavvk/postgressql/config"
	"github.com/AVVKavvk/postgressql/models"
	"github.com/jackc/pgx/v5"
)

var MyAppDBURl = "postgres://myuser:mypassword@localhost:5432/myapp_db"

func AddUserService(user *models.User) (map[string]interface{}, error) {
	dbRegistry, err := config.GetDBRegistry(context.Background(), MyAppDBURl)
	if err != nil {
		return nil, err
	}
	pool := dbRegistry.MyApp

	tableName := "users"
	isUsersTableExists := config.MyAppTables.Exists(tableName)

	if !isUsersTableExists {
		return map[string]interface{}{
			"result": nil,
		}, errors.New("failed to add user")
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (name, email, username, age) VALUES ($1, $2, $3, $4)`,
		tableName,
	)
	result, err := pool.Exec(context.Background(),
		query,
		user.Name, user.Email, user.UserName, user.Age,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"result": result,
	}, nil
}

func GetAllUsersService() ([]*models.User, error) {

	var users []*models.User

	dbRegistry, err := config.GetDBRegistry(context.Background(), MyAppDBURl)
	if err != nil {
		return nil, err
	}
	tableName := "users"
	isUsersTableExists := config.MyAppTables.Exists(tableName)

	if !isUsersTableExists {
		return []*models.User{}, errors.New("failed to add user")
	}

	pool := dbRegistry.MyApp

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	result, err := pool.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	for result.Next() {
		var user models.User
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.UserName, &user.Age, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}
func GetUserByIdService(id string) (*models.User, error) {
	var user *models.User = &models.User{}

	dbRegistry, err := config.GetDBRegistry(context.Background(), MyAppDBURl)
	if err != nil {
		return nil, err
	}
	tableName := "users"
	isUsersTableExists := config.MyAppTables.Exists(tableName)

	if !isUsersTableExists {
		return nil, errors.New("failed to add user")
	}

	pool := dbRegistry.MyApp

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

	result := pool.QueryRow(context.Background(), query, id)

	err = result.Scan(&user.ID, &user.Name, &user.Email, &user.UserName, &user.Age, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err // DB error
	}

	return user, nil

}

func UpdateUserService(id string, user *models.User) (*models.User, error) {

	dbRegistry, err := config.GetDBRegistry(context.Background(), MyAppDBURl)
	if err != nil {
		return nil, err
	}
	tableName := "users"
	isUsersTableExists := config.MyAppTables.Exists(tableName)

	if !isUsersTableExists {
		return nil, errors.New("failed to add user")
	}

	pool := dbRegistry.MyApp

	query := fmt.Sprintf("UPDATE %s SET name = $1, email = $2, username = $3, age = $4 WHERE id = $5", tableName)

	_, err = pool.Exec(context.Background(), query, user.Name, user.Email, user.UserName, user.Age, id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UserName:  user.UserName,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
	}, nil
}

func DeleteUserService(id string) error {

	dbRegistry, err := config.GetDBRegistry(context.Background(), MyAppDBURl)
	if err != nil {
		return err
	}
	tableName := "users"
	isUsersTableExists := config.MyAppTables.Exists(tableName)

	if !isUsersTableExists {
		return errors.New("failed to add user")
	}

	pool := dbRegistry.MyApp

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	_, err = pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
