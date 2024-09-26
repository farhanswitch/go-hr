package employee

import (
	"errors"
	"net/http"

	hashidutility "github.com/farhanswitch/go-hr/utilities/hashid"
)

type employeeService struct {
	repo employeeRepo
}

var service employeeService

func (es employeeService) getAllSrvc(param getAllEmployeeRequest) ([]employeeDataResponse, uint16, uint16, error) {
	listData, count, total, err := es.repo.GetAll(param)
	for i := range listData {
		listData[i].Encode()
	}
	return listData, count, total, err
}
func (es employeeService) getEmployeeDetails(employeeId uint16) (employeeDetailsResponse, error) {
	return es.repo.GetEmployeeDetails(employeeId)
}
func (es employeeService) createEmployee(data createEmployeeRequest) (int, map[string][1]string, error) {

	minSalary, maxSalary, err := es.repo.IsJobIdExists(data.DecodedJobId)
	if err != nil {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"jobId": {err.Error()},
		}, err
	}
	if data.Salary < minSalary || data.Salary > maxSalary {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"salary": {"Salary must be in range minimum salary and maximum salary for selected position."},
		}, errors.New("invalid salary")
	}
	if data.DecodedManagerId != 0 {
		err = es.repo.IsManagerIdExists(data.DecodedManagerId)
		if err != nil {
			return http.StatusUnprocessableEntity, map[string][1]string{
				"managerId": {err.Error()},
			}, err
		}

	}
	err = es.repo.IsDepartmentIDExists(data.DecodedDepartmentId)
	if err != nil {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"departmentID": {err.Error()},
		}, err
	}

	id, err := es.repo.AddEmployee(data)
	if err != nil {
		return http.StatusInternalServerError, map[string][1]string{
			"message": {err.Error()},
		}, err
	}
	hash := hashidutility.FactoryHashID()
	strUserId, err := hash.Encode([]int{int(id)})
	if err != nil {
		return http.StatusInternalServerError, map[string][1]string{
			"message": {err.Error()},
		}, err
	}
	return 0, map[string][1]string{
		"message":    {"Employee added successfully"},
		"employeeID": {strUserId},
	}, nil
}
func (es employeeService) updateEmployee(data updateEmployeeRequest) (int, map[string][1]string, error) {
	err := es.repo.IsEmployeeIDExists(data.DecodedEmployeeID)
	if err != nil {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"employeeId": {err.Error()},
		}, err
	}
	minSalary, maxSalary, err := es.repo.IsJobIdExists(data.DecodedJobId)
	if err != nil {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"jobId": {err.Error()},
		}, err
	}
	if data.Salary < minSalary || data.Salary > maxSalary {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"salary": {"Salary must be in range minimum salary and maximum salary for selected position."},
		}, errors.New("invalid salary")
	}
	if data.DecodedManagerId != 0 {
		err = es.repo.IsManagerIdExists(data.DecodedManagerId)
		if err != nil {
			return http.StatusUnprocessableEntity, map[string][1]string{
				"managerId": {err.Error()},
			}, err
		}

	}
	err = es.repo.IsDepartmentIDExists(data.DecodedDepartmentId)
	if err != nil {
		return http.StatusUnprocessableEntity, map[string][1]string{
			"departmentID": {err.Error()},
		}, err
	}

	err = es.repo.UpdateEmployee(data)
	if err != nil {
		return http.StatusInternalServerError, map[string][1]string{
			"message": {err.Error()},
		}, err
	}
	return 0, map[string][1]string{
		"message": {"Employee updated successfully"},
	}, nil
}
func factoryEmployeeService(repo employeeRepo) employeeService {
	if service == (employeeService{}) {
		service = employeeService{repo}
	}
	return service
}
