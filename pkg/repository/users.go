package repository

import (
	"context"
	"taskRumbler/pkg/models"
)

func (repo *PGRepo) CreateUser(user models.User) (id int, err error) {
	err = repo.pool.QueryRow(context.Background(), `INSERT INTO users (name, surname, patronymic, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id;`, user.Name, user.Surname, user.Patronymic, user.Email, user.Password).Scan(&id)
	return id, err
}

func (repo *PGRepo) GetUserById(id int) (user *models.User, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, name, surname, patronymic, email, password FROM users WHERE id =$1`, id).Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Email,
		&user.Password,
	)
	return user, err
}

func (repo *PGRepo) GetUsers() ([]models.User, error) {
	var usersData []models.User
	row, err := repo.pool.Query(context.Background(), `SELECT id, name, surname, patronymic, email, password FROM users`)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		var user models.User
		err = row.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		usersData = append(usersData, user)
	}
	return usersData, nil
}

func (repo *PGRepo) GetUserByEmail(email string) (models.User, error) {
	if email == "" {
		return models.User{}, nil
	}
	var user models.User
	err := repo.pool.QueryRow(context.Background(), `SELECT id, name, surname, patronymic, email, password FROM users WHERE email =$1`, email).Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Email,
		&user.Password,
	)
	return user, err
}
