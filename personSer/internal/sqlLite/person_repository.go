package sqlLite

import (
	"PersonService/internal/app/model"
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

type SqliteSQLPersonRepository struct {
	connection *sql.DB
}

func NewSqlLiteSQLPersonRepository(connection *sql.DB) model.PersonRepository {
	return &SqliteSQLPersonRepository{connection}
}

func ErrorHandler(rows *sql.Rows) {
	errRow := rows.Close()
	if errRow != nil {
		logrus.Error(errRow)
	}
}

func (p *SqliteSQLPersonRepository) Get(context context.Context) (res []model.Person, err error) {
	// SQL-запрос для выборки всех персон, отсортированных по ID
	query := `SELECT *
			  FROM Person 
			  ORDER BY id`
	// Выполняем запрос с контекстом
	rows, err := p.connection.QueryContext(context, query)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Закрываем строки результата после выполнения запроса
	defer ErrorHandler(rows)

	res = make([]model.Person, 0)
	// Итерируемся по строкам результата и сканируем данные в структуру Person
	for rows.Next() {
		temp := model.Person{}
		err = rows.Scan(
			&temp.ID,
			&temp.Email,
			&temp.Phone,
			&temp.FirstName,
			&temp.LastName,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (p *SqliteSQLPersonRepository) GetByID(context context.Context, id int64) (res model.Person, err error) {
	// SQL-запрос для выборки персоны по ID
	query := `SELECT *
			  FROM Person 
			  WHERE id = ?`
	// Выполняем запрос с контекстом и параметром ID
	rows, err := p.connection.QueryContext(context, query, id)

	if err != nil {
		logrus.Error(err)
		return res, err
	}

	// Закрываем строки результата после выполнения запроса
	defer ErrorHandler(rows)

	if rows.Next() {
		// Если есть результат, сканируем его в структуру Person
		err = rows.Scan(
			&res.ID,
			&res.Email,
			&res.Phone,
			&res.FirstName,
			&res.LastName,
		)
	} else {
		// Если нет результата, логируем ошибку и возвращаем ошибку "person not found"
		logrus.Error(errors.New("person not found"))
		return res, errors.New("person not found")
	}
	return res, nil
}

func (p *SqliteSQLPersonRepository) Add(context context.Context, person model.Person) (id int64, err error) {
	// SQL-запрос для добавления новой персоны
	query := `INSERT INTO Person (email, phone, firstname, lastname) 
			  VALUES (?, ?, ?, ?)`

	// Выполняем запрос с контекстом и параметрами новой персоны
	res, err := p.connection.ExecContext(context, query, person.Email, person.Phone, person.FirstName, person.LastName)

	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	// Проверяем количество затронутых строк
	err = rowsAffectedCheck(res)
	if err != nil {
		return 0, err
	}

	// Получаем ID новой персоны
	lastID, err := res.LastInsertId()

	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	id = int64(lastID)

	return id, nil
}

func (p *SqliteSQLPersonRepository) Update(context context.Context, person model.Person) (err error) {
	// SQL-запрос для обновления персоны по ID
	query := `UPDATE Person SET email=?, phone=?, firstname=?, lastname=? WHERE id = ?`

	// Выполняем запрос с контекстом и параметрами обновления персоны
	res, err := p.connection.ExecContext(context, query, person.Email, person.Phone, person.FirstName, person.LastName, person.ID)

	if err != nil {
		logrus.Error(err)
		return err
	}

	// Проверяем количество затронутых строк
	err = rowsAffectedCheck(res)
	if err != nil {
		return err
	}

	return nil
}

func (p *SqliteSQLPersonRepository) Delete(context context.Context, id int64) (err error) {
	// SQL-запрос для удаления персоны по ID
	query := `DELETE FROM Person WHERE id = ?`
	res, err := p.connection.ExecContext(context, query, id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// Проверяем количество затронутых строк
	err = rowsAffectedCheck(res)
	if err != nil {
		return err
	}

	return nil
}

func rowsAffectedCheck(res sql.Result) (err error) {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return err
	}
	if rowsAffected != 1 {
		logrus.Error(errors.New("more than one row affected"))
		return err
	}
	return nil
}
