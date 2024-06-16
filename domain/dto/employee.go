package dto

import "time"

type CreateEmployeeResponse struct {
	Id int `json:"id"`
}

type DeleteEmployeeRequest struct {
	EmployeeID int `json:"employee_id" uri:"employee_id" binding:"required"`
}

type GetEmployee struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type GetEmployeeByIDRequest struct {
	EmployeeID int `json:"employee_id" uri:"employee_id" binding:"required"`
}

type EmployeeCreateRequest struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

// Employee Represents the fields from the Employee Database
type Employee struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	Salary    float64   `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateEmployeeBodyRequest struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

type UpdateEmployeeRequest struct {
	EmployeeID int `json:"employee_id" uri:"employee_id" binding:"required"`
}
