package usecase

import (
	"database/sql"
	"employee-management/api/repository/sqlboiler"
	"employee-management/domain/dto"
	"employee-management/domain/interfaces"
	"employee-management/utils/convert"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type employeeUsecase struct {
	db *sql.DB
}

func NewEmployeeUsecase(db *sql.DB) interfaces.EmployeeUsecase {
	return &employeeUsecase{
		db: db,
	}
}

func (uc *employeeUsecase) GetEmployeeById(ctx *gin.Context, employeeID int) (*dto.Employee, error) {
	boil.DebugMode = true
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	employee, err := sqlboiler.FindEmployee(ctx, tx, employeeID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return convert.ToEmployeeDTO(employee), nil

}

func (uc *employeeUsecase) GetAllEmployee(ctx *gin.Context, limit int, offset int) ([]*dto.Employee, error) {
	boil.DebugMode = true
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	employee, err := sqlboiler.Employees(qm.Limit(limit), qm.Offset(offset)).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return convert.ToEmployeeSliceDTO(employee), nil
}

func (uc *employeeUsecase) CreateEmployee(ctx *gin.Context, request *dto.EmployeeCreateRequest) (dto.CreateEmployeeResponse, error) {
	boil.DebugMode = true
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return dto.CreateEmployeeResponse{}, err
	}
	defer tx.Rollback()

	employee := sqlboiler.Employee{
		Name:      request.Name,
		Position:  request.Position,
		Salary:    request.Salary,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = employee.Insert(ctx, uc.db, boil.Infer())
	if err != nil {
		return dto.CreateEmployeeResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return dto.CreateEmployeeResponse{}, err
	}

	return dto.CreateEmployeeResponse{
		Id: employee.ID,
	}, nil
}

func (uc *employeeUsecase) UpdateEmployee(ctx *gin.Context, employeeID int, request *dto.UpdateEmployeeBodyRequest) (*dto.Employee, error) {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	employee := sqlboiler.M{
		"updated_at": time.Now(),
	}

	if len(request.Name) != 0 {
		employee["name"] = request.Name
	}

	if len(request.Position) != 0 {
		employee["position"] = request.Position
	}

	if request.Salary != 0 {
		employee["salary"] = request.Salary
	}

	_, err = sqlboiler.Employees(qm.Where("id=?", employeeID)).UpdateAll(ctx, uc.db, employee)
	if err != nil {
		return &dto.Employee{}, err
	}

	emp, err := sqlboiler.Employees(qm.Where("id=?", employeeID)).One(ctx, uc.db)
	if err != nil {
		return nil, fmt.Errorf("could not find employee with id %d: %w", employeeID, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Map the retrieved model to your Employee struct
	employeedata := &dto.Employee{
		ID:        emp.ID,
		Name:      emp.Name,
		Position:  emp.Position,
		Salary:    emp.Salary,
		CreatedAt: emp.CreatedAt,
		UpdatedAt: emp.UpdatedAt,
	}

	return employeedata, nil
}

func (uc *employeeUsecase) DeleteEmployee(ctx *gin.Context, employeeID int) error {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	employee, err := sqlboiler.FindEmployee(ctx, uc.db, employeeID)
	if err != nil {
		return err
	}

	_, err = sqlboiler.Employees(qm.Where("id = ?", employee.ID)).DeleteAll(ctx, uc.db)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
