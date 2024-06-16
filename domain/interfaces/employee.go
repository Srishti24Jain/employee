package interfaces

import (
	"employee-management/domain/dto"

	"github.com/gin-gonic/gin"
)

type EmployeeUsecase interface {
	GetEmployeeById(ctx *gin.Context, employeeID int) (*dto.Employee, error)
	GetAllEmployee(ctx *gin.Context, limit int, offset int) ([]*dto.Employee, error)
	CreateEmployee(ctx *gin.Context, request *dto.EmployeeCreateRequest) (dto.CreateEmployeeResponse, error)
	UpdateEmployee(ctx *gin.Context, employeeID int, requestBody *dto.UpdateEmployeeBodyRequest) (*dto.Employee, error)
	DeleteEmployee(ctx *gin.Context, employeeID int) error
}
