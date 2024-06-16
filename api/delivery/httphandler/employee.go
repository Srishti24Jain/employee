package httphandler

import (
	"employee-management/domain/dto"
	"employee-management/domain/interfaces"
	"employee-management/utils/httputil"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type employeeHandler struct {
	employeeUsecase interfaces.EmployeeUsecase
}

func NewEmployeeHandler(e *gin.Engine, a interfaces.EmployeeUsecase) {
	handler := employeeHandler{employeeUsecase: a}
	fmt.Println("handler", handler)
	e.GET("api/employee/:employee_id", handler.GetEmployeeByIdHandler)
	e.GET("api/list_employee", handler.GetEmployeeHandler)
	e.POST("api/add-employee", handler.CreateEmployeeHandler)
	e.PUT("api/employee/:employee_id", handler.UpdateEmployeeHandler)
	e.DELETE("api/employee/:employee_id", handler.DeleteEmployeeHandler)
}

func (s *employeeHandler) GetEmployeeByIdHandler(ctx *gin.Context) {
	var (
		startTime = time.Now()
		httpError *httputil.StandardError
	)
	defer func() {
		if httpError != nil {
			errCode, _ := strconv.Atoi(httpError.Code)
			httputil.WriteErrorResponse(ctx.Writer, errCode, []httputil.StandardError{*httpError})
		}
	}()

	req := new(dto.GetEmployeeByIDRequest)

	if err := ctx.ShouldBindUri(req); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}

	employee, err := s.employeeUsecase.GetEmployeeById(ctx, req.EmployeeID)
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	data, err := json.Marshal(httputil.StandardEnvelope{
		Data: employee,
		Status: &httputil.StandardStatus{
			Message:   http.StatusText(http.StatusOK),
			ErrorCode: 0,
		},
		Header: &httputil.StandardHeader{
			TotalData:   1,
			ProcessTime: time.Since(startTime).Seconds(),
		},
	})
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}
	_, _ = httputil.WriteJSONResponse(ctx.Writer, data, http.StatusOK)
	return
}

func (s *employeeHandler) GetEmployeeHandler(ctx *gin.Context) {
	var (
		startTime = time.Now()
		httpError *httputil.StandardError
	)
	defer func() {
		if httpError != nil {
			errCode, _ := strconv.Atoi(httpError.Code)
			httputil.WriteErrorResponse(ctx.Writer, errCode, []httputil.StandardError{*httpError})
		}
	}()

	req := new(dto.GetEmployee)
	if err := ctx.ShouldBindQuery(&req); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}

	if req.PageSize == 0 {
		req.PageSize = 100
	}
	offset := (req.Page - 1) * req.PageSize
	limit := req.PageSize
	fmt.Println("offset", offset)

	employee, err := s.employeeUsecase.GetAllEmployee(ctx, limit, offset)
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	data, err := json.Marshal(httputil.StandardEnvelope{
		Data: employee,
		Status: &httputil.StandardStatus{
			Message:   http.StatusText(http.StatusOK),
			ErrorCode: 0,
		},
		Header: &httputil.StandardHeader{
			TotalData:   1,
			ProcessTime: time.Since(startTime).Seconds(),
		},
	})
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	_, _ = httputil.WriteJSONResponse(ctx.Writer, data, http.StatusOK)
	return
}

func (s *employeeHandler) CreateEmployeeHandler(ctx *gin.Context) {
	var (
		startTime = time.Now()
		httpError *httputil.StandardError
	)
	defer func() {
		if httpError != nil {
			errCode, _ := strconv.Atoi(httpError.Code)
			httputil.WriteErrorResponse(ctx.Writer, errCode, []httputil.StandardError{*httpError})
		}
	}()

	req := new(dto.EmployeeCreateRequest)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}

	resp, err := s.employeeUsecase.CreateEmployee(ctx, req)
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	data, err := json.Marshal(httputil.StandardEnvelope{
		Data: resp,
		Status: &httputil.StandardStatus{
			Message:   http.StatusText(http.StatusOK),
			ErrorCode: 0,
		},
		Header: &httputil.StandardHeader{
			TotalData:   1,
			ProcessTime: time.Since(startTime).Seconds(),
		},
	})
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}
	_, _ = httputil.WriteJSONResponse(ctx.Writer, data, http.StatusCreated)
	return
}

func (s *employeeHandler) UpdateEmployeeHandler(ctx *gin.Context) {
	var (
		startTime = time.Now()
		httpError *httputil.StandardError
	)
	defer func() {
		if httpError != nil {
			errCode, _ := strconv.Atoi(httpError.Code)
			httputil.WriteErrorResponse(ctx.Writer, errCode, []httputil.StandardError{*httpError})
		}
	}()

	req := new(dto.UpdateEmployeeRequest)
	if err := ctx.BindUri(req); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}

	reqBody := new(dto.UpdateEmployeeBodyRequest)
	if err := ctx.Bind(&reqBody); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}

	resp, err := s.employeeUsecase.UpdateEmployee(ctx, req.EmployeeID, reqBody)
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	data, err := json.Marshal(httputil.StandardEnvelope{
		Data: resp,
		Status: &httputil.StandardStatus{
			Message:   http.StatusText(http.StatusOK),
			ErrorCode: 0,
		},
		Header: &httputil.StandardHeader{
			TotalData:   0,
			ProcessTime: time.Since(startTime).Seconds(),
		},
	})
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}
	_, _ = httputil.WriteJSONResponse(ctx.Writer, data, http.StatusOK)
	return
}

func (s *employeeHandler) DeleteEmployeeHandler(ctx *gin.Context) {
	var (
		startTime = time.Now()
		httpError *httputil.StandardError
	)
	defer func() {
		if httpError != nil {
			errCode, _ := strconv.Atoi(httpError.Code)
			httputil.WriteErrorResponse(ctx.Writer, errCode, []httputil.StandardError{*httpError})
		}
	}()
	req := new(dto.DeleteEmployeeRequest)
	if err := ctx.ShouldBindUri(req); err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusBadRequest),
			Title:  http.StatusText(http.StatusBadRequest),
			Detail: err.Error(),
		}
		return
	}
	err := s.employeeUsecase.DeleteEmployee(ctx, req.EmployeeID)
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}

	data, err := json.Marshal(httputil.StandardEnvelope{
		Data: "Employee Deleted Successfully",
		Status: &httputil.StandardStatus{
			Message:   http.StatusText(http.StatusNoContent),
			ErrorCode: http.StatusNoContent,
		},
		Header: &httputil.StandardHeader{
			TotalData:   1,
			ProcessTime: time.Since(startTime).Seconds(),
		},
	})
	if err != nil {
		httpError = &httputil.StandardError{
			Code:   strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: err.Error(),
		}
		return
	}
	_, _ = httputil.WriteJSONResponse(ctx.Writer, data, http.StatusOK)
	return
}
