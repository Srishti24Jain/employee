package usecase

import (
	"employee-management/domain/dto"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployeeById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	employeeID := 1
	expectedEmployee := &dto.Employee{
		ID:        employeeID,
		Name:      "John Doe",
		Position:  "Developer",
		Salary:    60000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "position", "salary", "created_at", "updated_at"}).
		AddRow(employeeID, "John Doe", "Developer", 60000, time.Now(), time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from "employee" where "id"=$1`)).WithArgs(employeeID).WillReturnRows(rows)
	mock.ExpectCommit()

	uc := NewEmployeeUsecase(db)
	ctx := &gin.Context{}

	employee, err := uc.GetEmployeeById(ctx, employeeID)

	assert.NoError(t, err)
	assert.NotNil(t, employee)
	assert.Equal(t, expectedEmployee.Name, employee.Name)
	assert.Equal(t, expectedEmployee.Position, employee.Position)
	assert.Equal(t, expectedEmployee.Salary, employee.Salary)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedEmployees := []*dto.Employee{
		{
			ID:        1,
			Name:      "John Doe",
			Position:  "Developer",
			Salary:    60000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Jane Smith",
			Position:  "Manager",
			Salary:    80000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "position", "salary", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "Developer", 60000, time.Now(), time.Now()).
		AddRow(2, "Jane Smith", "Manager", 80000, time.Now(), time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "employee".* FROM "employee" LIMIT 10 OFFSET 5`)).WillReturnRows(rows)
	mock.ExpectCommit()

	uc := NewEmployeeUsecase(db)
	ctx := &gin.Context{}

	employees, err := uc.GetAllEmployee(ctx, 10, 5)

	assert.NoError(t, err)
	assert.Len(t, employees, 2)
	assert.Equal(t, expectedEmployees[0].Name, employees[0].Name)
	assert.Equal(t, expectedEmployees[1].Name, employees[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id int
	id = 1
	request := &dto.EmployeeCreateRequest{
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "employee" ("name","position","salary","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(request.Name, request.Position, request.Salary, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	uc := NewEmployeeUsecase(db)
	ctx := &gin.Context{}

	response, err := uc.CreateEmployee(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, id, response.Id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	employeeID := 1
	request := &dto.UpdateEmployeeBodyRequest{
		Name:     "John Doe Updated",
		Position: "Senior Developer",
		Salary:   70000,
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "employee" SET "name" = $1, "position" = $2, "salary" = $3, "updated_at" = $4 WHERE (id=$5)`)).
		WithArgs(request.Name, request.Position, request.Salary, sqlmock.AnyArg(), employeeID).
		WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock the select query after update
	rows := sqlmock.NewRows([]string{"id", "name", "position", "salary", "created_at", "updated_at"}).
		AddRow(employeeID, request.Name, request.Position, request.Salary, time.Now(), time.Now())
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "employee".* FROM "employee" WHERE (id=$1) LIMIT 1`)).WithArgs(employeeID).WillReturnRows(rows)

	mock.ExpectCommit()

	uc := NewEmployeeUsecase(db)
	ctx := &gin.Context{}

	employee, err := uc.UpdateEmployee(ctx, employeeID, request)

	assert.NoError(t, err)
	assert.Equal(t, request.Name, employee.Name)
	assert.Equal(t, request.Position, employee.Position)
	assert.Equal(t, request.Salary, employee.Salary)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	employeeID := 1

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from "employee" where "id"=$1`)).
		WithArgs(employeeID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary", "created_at", "updated_at"}).
			AddRow(employeeID, "John Doe", "Developer", 60000, time.Now(), time.Now()))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "employee" WHERE (id = $1)`)).
		WithArgs(employeeID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uc := NewEmployeeUsecase(db)
	ctx := &gin.Context{}

	err = uc.DeleteEmployee(ctx, employeeID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
