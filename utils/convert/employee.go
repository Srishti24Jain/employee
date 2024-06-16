package convert

import (
	"employee-management/api/repository/sqlboiler"
	"employee-management/domain/dto"
)

func ToEmployeeDTO(employee *sqlboiler.Employee) *dto.Employee {
	e := &dto.Employee{
		ID:        employee.ID,
		Name:      employee.Name,
		Position:  employee.Position,
		Salary:    employee.Salary,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}

	return e
}

func ToEmployeeSliceDTO(employeeSlice sqlboiler.EmployeeSlice) []*dto.Employee {
	allemployee := make([]*dto.Employee, 0, len(employeeSlice))
	for _, employee := range employeeSlice {
		allemployee = append(allemployee, ToEmployeeDTO(employee))
	}

	return allemployee
}
