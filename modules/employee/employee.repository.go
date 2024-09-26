package employee

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	hashidutility "github.com/farhanswitch/go-hr/utilities/hashid"
)

type employeeRepo struct {
	db *sql.DB
}

var repo employeeRepo

func (er employeeRepo) GetAll(param getAllEmployeeRequest) ([]employeeDataResponse, uint16, uint16, error) {
	results, err := er.db.Query("CALL sp_get_all_employee(?,?,?,?,?)", param.Search, ((param.PaginationPage - 1) * param.PaginationRows), param.PaginationRows, param.SortField, param.SortOrder)
	if err != nil {
		return []employeeDataResponse{}, 0, 0, err
	}
	var listEmployee []employeeDataResponse = []employeeDataResponse{}
	//* Handle collections of employee
	for results.Next() {
		var data employeeDataResponse
		var strTime string
		var nullableManagerName sql.NullString
		var nullablePhoneNumber sql.NullString
		err := results.Scan(&data.EmployeeID, &data.FirstName, &data.LastName, &data.Email, &nullablePhoneNumber, &strTime, &data.JobTitle, &data.Salary, &nullableManagerName, &data.DepartmentName)
		if err != nil {
			return []employeeDataResponse{}, 0, 0, err
		}
		//* Handle Nullable Manager Id
		if nullableManagerName.Valid {
			data.ManagerName = nullableManagerName.String
		} else {
			data.ManagerName = ""
		}
		//* Handle Nullable Phone Number
		if nullablePhoneNumber.Valid {
			data.PhoneNumber = nullablePhoneNumber.String
		} else {
			data.PhoneNumber = ""
		}
		//* Handle Parsing string of date to time.Time
		hireDate, err := time.Parse("2006-01-02", strTime)
		if err != nil {
			return []employeeDataResponse{}, 0, 0, err
		}
		data.HireDate = hireDate
		listEmployee = append(listEmployee, data)
	}
	var listCount []uint16 = []uint16{}
	for results.NextResultSet() {
		if results.Next() {
			var count uint16
			err := results.Scan(&count)
			if err != nil {
				return []employeeDataResponse{}, 0, 0, err
			}
			listCount = append(listCount, count)
		}
	}

	countCurrent := listCount[0]
	countTotal := listCount[1]
	results.Close()
	return listEmployee, countCurrent, countTotal, nil
}

func (er employeeRepo) GetEmployeeDetails(employeeId uint16) (employeeDetailsResponse, error) {
	log.Printf("EmployeeId: %d\n", employeeId)
	employee := employeeResponse{}
	manager := employeeManagerResponse{}
	department := employeeDepartmentResponse{}
	job := employeeJobResponse{}
	employeeDetails := employeeDetailsResponse{}
	var strDate string
	err := er.db.QueryRow("CALL sp_get_employee_details(?)", employeeId).Scan(&employee.EmployeeId, &employee.FirstName, &employee.LastName, &employee.Email, &employee.PhoneNumber, &strDate, &employee.Salary, &department.DepartmentId, &department.DepartmentName, &department.DepartmentLocationId, &job.JobId, &job.JobTitle, &manager.ManagerId, &manager.ManagerFirstName, &manager.ManagerLastName)

	if err != nil {
		return employeeDetails, err
	}
	//* Handle Parsing string of date to time.Time
	hireDate, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		return employeeDetailsResponse{}, err
	}
	employee.HireDate = hireDate
	hash := hashidutility.FactoryHashID()
	strEmployeeId, _ := hash.Encode([]int{int(employee.EmployeeId)})
	strDepartmentId, _ := hash.Encode([]int{int(department.DepartmentId)})
	strJobId, _ := hash.Encode([]int{int(job.JobId)})
	strManagerId, _ := hash.Encode([]int{int(manager.ManagerId)})
	strDepartmentLocationId, _ := hash.Encode([]int{int(department.DepartmentLocationId)})
	employee.HashedEmployeeId = strEmployeeId
	department.HashedDepartmentId = strDepartmentId
	department.HashedDepartmentLocationId = strDepartmentLocationId
	job.HashedJobId = strJobId
	manager.HashedManagerID = strManagerId

	employeeDetails = employeeDetailsResponse{
		EmployeeDetails:   employee,
		DepartmentDetails: department,
		JobDetails:        job,
		ManagerDetails:    manager,
	}
	return employeeDetails, nil
}
func (er employeeRepo) IsJobIdExists(jobId uint16) (float32, float32, error) {
	var minSalary, maxSalary float32
	if err := er.db.QueryRow("SELECT min_salary, max_salary FROM jobs WHERE job_id = ? LIMIT 1;", jobId).Scan(&minSalary, &maxSalary); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("there is no job with id %d", jobId)
		} else {
			return 0, 0, err
		}
	}
	return minSalary, maxSalary, nil
}
func (er employeeRepo) IsManagerIdExists(managerId uint16) error {
	var id uint16
	if err := er.db.QueryRow("SELECT employee_id FROM employees WHERE employee_id = ? LIMIT 1;", managerId).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("there is no manager with id %d", managerId)
		} else {
			return err
		}
	}
	if id == 0 {
		return errors.New("invalid manager id")
	}
	return nil
}
func (er employeeRepo) IsDepartmentIDExists(departmentID uint16) error {
	var id uint16
	if err := er.db.QueryRow("SELECT department_id FROM departments WHERE department_id = ? LIMIT 1;", departmentID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("there is no department wih id %d", departmentID)
		} else {
			return err
		}
	}
	if id <= 0 {
		return errors.New("invalid department ID")
	}
	return nil
}
func (er employeeRepo) IsEmployeeIDExists(employeeID uint16) error {
	var id uint16
	if err := er.db.QueryRow("SELECT employee_id FROM employees WHERE employee_id = ? LIMIT 1;", employeeID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("there is no employee wih id %d", employeeID)
		} else {
			return err
		}
	}
	if id <= 0 {
		return errors.New("invalid employee ID")
	}
	return nil
}
func (er employeeRepo) AddEmployee(data createEmployeeRequest) (uint16, error) {
	var id uint16
	if err := er.db.QueryRow("CALL sp_create_employee(?,?,?,?,?,?,?,?,?)", data.FirstName, data.LastName, data.EmailAddress, data.PhoneNumber, data.HireDate, data.DecodedJobId, data.DecodedManagerId, data.DecodedDepartmentId, data.Salary).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (er employeeRepo) UpdateEmployee(data updateEmployeeRequest) error {
	_, err := er.db.Exec("CALL sp_update_employee(?,?,?,?,?,?,?,?,?,?);", data.DecodedEmployeeID, data.FirstName, data.LastName, data.EmailAddress, data.PhoneNumber, data.HireDate, data.DecodedJobId, data.DecodedManagerId, data.DecodedDepartmentId, data.Salary)
	return err
}
func factoryEmployeeRepository(db *sql.DB) employeeRepo {
	if repo == (employeeRepo{}) {
		repo = employeeRepo{db}
	}
	return repo
}
