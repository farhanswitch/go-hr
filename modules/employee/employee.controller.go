package employee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	responsehelper "github.com/farhanswitch/go-hr/helpers/response"
	errorutility "github.com/farhanswitch/go-hr/utilities/errors"
	hashidutility "github.com/farhanswitch/go-hr/utilities/hashid"
)

type employeeController struct {
	service  employeeService
	validate *validator.Validate
}

var controller employeeController

func (ec employeeController) getAllCtrl(w http.ResponseWriter, r *http.Request) {
	sortOrder := r.URL.Query().Get("sortOrder")
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sortField")
	if sortField == "" {
		sortField = "first_name"
	}
	w.Header().Set("Content-Type", "application/json")

	var errorObject map[string][1]string = make(map[string][1]string)

	paginationPage, err := strconv.Atoi(r.URL.Query().Get("paginationPage"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObject = map[string][1]string{
			"paginationPage": {"paginationPage must be a number"},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors": %s}`, errString)))
		return
	}
	paginationRows, err := strconv.Atoi(r.URL.Query().Get("paginationRows"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObject = map[string][1]string{
			"paginationRows": {"paginationRows must be a number"},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}
	if paginationPage < 1 {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"paginationPage": {"paginationPage must be greater than or equal to 1."},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors": %s}`, errString)))
		return
	}

	requestParams := getAllEmployeeRequest{
		Search:         search,
		PaginationPage: uint16(paginationPage),
		PaginationRows: uint16(paginationRows),
		SortField:      sortField,
		SortOrder:      sortOrder,
	}

	err = ec.validate.Struct(requestParams)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorRes := errorutility.ParseError(err)
		errString, _ := json.Marshal(errorRes)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return

	}

	res, countCurrent, countTotal, err := ec.service.getAllSrvc(requestParams)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal Server Error."}`))
		return
	}

	jsonData, _ := json.Marshal(res)
	w.Write([]byte(fmt.Sprintf(`{"data":%s,"count":%d,"total":%d}`, jsonData, countCurrent, countTotal)))
}
func (ec employeeController) getItemCtrl(w http.ResponseWriter, r *http.Request) {
	employeeId := chi.URLParam(r, "employeeId")
	if employeeId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"employeeId": {"employeeId is required."},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return

	}
	hash := hashidutility.FactoryHashID()
	numEmployeeId, err := hash.DecodeWithError(employeeId)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		err = errorutility.DetectOtherError(err)
		errObject := map[string][1]string{
			"employeeId": {err.Error()},
		}
		errString, _ := json.Marshal(errObject)
		w.Write([]byte(fmt.Sprintf(`{"errors": %s}`, errString)))
		return
	}
	employeeDetails, err := ec.service.getEmployeeDetails(uint16(numEmployeeId[0]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getEmployeeDetails with id %s %d.\nError: %s", employeeId, numEmployeeId, err.Error())
		w.Write([]byte(`{"errors":"Internal Server Error."}`))
		return
	}
	jsonData, _ := json.Marshal(employeeDetails)
	w.Write([]byte(fmt.Sprintf(`{"data":%s}`, jsonData)))
}

func (ec employeeController) createCtrl(w http.ResponseWriter, r *http.Request) {
	var reqBody createEmployeeRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		responsehelper.SimpleValidationError("createEmployeeCtrl", err, w)
		return
	}
	hash := hashidutility.FactoryHashID()
	if reqBody.JobId == "" {
		responsehelper.SimpleValidationError("createEmployeeCtrl", fmt.Errorf("jobId cannot be empty"), w)
		return
	} else if reqBody.ManagerId == "" {
		responsehelper.SimpleValidationError("createEmployeeCtrl", fmt.Errorf("managerId cannot be empty"), w)
		return
	} else if reqBody.DepartmentId == "" {
		responsehelper.SimpleValidationError("createEmployeeCtrl", fmt.Errorf("departmentId cannot be empty"), w)
		return
	}
	jobId, err := hash.DecodeWithError(reqBody.JobId)
	if err != nil {
		responsehelper.SimpleValidationError("createEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedJobId = uint16(jobId[0])

	managerId, err := hash.DecodeWithError(reqBody.ManagerId)
	if err != nil {
		responsehelper.SimpleValidationError("createEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedManagerId = uint16(managerId[0])

	departmentId, err := hash.DecodeWithError(reqBody.DepartmentId)
	if err != nil {
		responsehelper.SimpleValidationError("createEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedDepartmentId = uint16(departmentId[0])
	err = ec.validate.Struct(reqBody)
	if err != nil {
		errString, _ := json.Marshal(errorutility.ParseError(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}

	layout := "2006-01-02"
	parsedHireDate, err := time.Parse(layout, reqBody.HireDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"hireDate": {"Invalid date."},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}
	today := time.Now()
	isValidDate := today.After(parsedHireDate)
	if !isValidDate {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"hireDate": {"Hire Date must be before today"},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}
	errorCode, data, err := ec.service.createEmployee(reqBody)

	if err != nil {
		log.Printf("Error createCtrl.\nError: %s\n", err.Error())

		if errorCode == http.StatusUnprocessableEntity {
			errJSON, _ := json.Marshal(data)
			errString := fmt.Sprintf(`{"errors": %s}`, errJSON)
			w.WriteHeader(errorCode)
			w.Write([]byte(errString))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors": "Internal Server Error."}`))
		return
	}
	jsonData, _ := json.Marshal(data)
	w.Write([]byte(fmt.Sprintf(`{"data": %s}`, jsonData)))
}
func (ec employeeController) updateCtrl(w http.ResponseWriter, r *http.Request) {
	var reqBody updateEmployeeRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", err, w)
		return
	}
	hash := hashidutility.FactoryHashID()
	if reqBody.EmployeeID == "" {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", fmt.Errorf("employeeId cannot be empty"), w)
		return
	} else if reqBody.JobId == "" {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", fmt.Errorf("jobId cannot be empty"), w)
		return
	} else if reqBody.ManagerId == "" {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", fmt.Errorf("managerId cannot be empty"), w)
		return
	} else if reqBody.DepartmentId == "" {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", fmt.Errorf("departmentId cannot be empty"), w)
		return
	}
	employeeId, err := hash.DecodeWithError(reqBody.EmployeeID)
	if err != nil {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedEmployeeID = uint16(employeeId[0])
	jobId, err := hash.DecodeWithError(reqBody.JobId)
	if err != nil {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedJobId = uint16(jobId[0])

	managerId, err := hash.DecodeWithError(reqBody.ManagerId)
	if err != nil {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedManagerId = uint16(managerId[0])

	departmentId, err := hash.DecodeWithError(reqBody.DepartmentId)
	if err != nil {
		responsehelper.SimpleValidationError("updateEmployeeCtrl", err, w)
		return
	}
	reqBody.DecodedDepartmentId = uint16(departmentId[0])
	err = ec.validate.Struct(reqBody)
	if err != nil {
		errString, _ := json.Marshal(errorutility.ParseError(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}

	layout := "2006-01-02"
	parsedHireDate, err := time.Parse(layout, reqBody.HireDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"hireDate": {"Invalid date."},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}
	today := time.Now()
	isValidDate := today.After(parsedHireDate)
	if !isValidDate {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := map[string][1]string{
			"hireDate": {"Hire Date must be before today"},
		}
		errString, _ := json.Marshal(errorObject)
		w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errString)))
		return
	}
	errorCode, data, err := ec.service.updateEmployee(reqBody)

	if err != nil {
		log.Printf("Error updateCtrl.\nError: %s\n", err.Error())

		if errorCode == http.StatusUnprocessableEntity {
			errJSON, _ := json.Marshal(data)
			errString := fmt.Sprintf(`{"errors": %s}`, errJSON)
			w.WriteHeader(errorCode)
			w.Write([]byte(errString))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors": "Internal Server Error."}`))
		return
	}
	jsonData, _ := json.Marshal(data)
	w.Write([]byte(fmt.Sprintf(`{"data": %s}`, jsonData)))
}

func factoryEmployeeController(service employeeService, validate *validator.Validate) employeeController {
	if controller == (employeeController{}) {
		controller = employeeController{service, validate}
	}
	return controller
}
