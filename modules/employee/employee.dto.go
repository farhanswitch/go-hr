package employee

import (
	"time"

	hashidutility "github.com/farhanswitch/go-hr/utilities/hashid"
)

type getAllEmployeeRequest struct {
	PaginationPage uint16 `validate:"required,gte=1"`
	PaginationRows uint16 `validate:"required,lte=1000,gte=1"`
	SortField      string `validate:"required,oneof=EmployeeID Salary FirstName LastName HireDate Email PhoneNumber JobTitle DepartmentName ManagerName"`
	SortOrder      string `validate:"required,oneof=asc desc ASC DESC"`
	Search         string
}

type employeeResponse struct {
	JobId            uint16 `json:"-"`
	ManagerId        uint16 `json:"-"`
	DepartmentId     uint16 `json:"-"`
	EmployeeId       uint16 `json:"-"`
	HashedEmployeeId string `json:"employeeId"`
	Salary           float32
	HireDate         time.Time
	FirstName        string
	LastName         string
	Email            string
	PhoneNumber      string
}
type employeeDataResponse struct {
	EmployeeID       uint16 `json:"-"`
	HashedEmployeeID string `json:"EmployeeID"`
	Salary           float32
	HireDate         time.Time
	FirstName        string
	LastName         string
	Email            string
	PhoneNumber      string
	JobTitle         string
	DepartmentName   string
	ManagerName      string
}

func (er *employeeDataResponse) Encode() {
	er.HashedEmployeeID, _ = hashidutility.FactoryHashID().Encode([]int{int(er.EmployeeID)})
}

type employeeManagerResponse struct {
	ManagerId        uint16 `json:"-"`
	HashedManagerID  string `json:"managerId"`
	ManagerFirstName string
	ManagerLastName  string
}
type employeeDepartmentResponse struct {
	DepartmentId uint16 `json:"-"`

	DepartmentLocationId       uint16 `json:"-"`
	HashedDepartmentLocationId string `json:"locationId"`
	HashedDepartmentId         string `json:"departmentId"`
	DepartmentName             string
}
type employeeJobResponse struct {
	JobId       uint16 `json:"-"`
	HashedJobId string `json:"jobId"`
	JobTitle    string
}
type employeeDetailsResponse struct {
	EmployeeDetails   employeeResponse
	ManagerDetails    employeeManagerResponse
	DepartmentDetails employeeDepartmentResponse
	JobDetails        employeeJobResponse
}
type createEmployeeRequest struct {
	DecodedJobId        uint16  `validate:"required,gte=1"`
	DecodedManagerId    uint16  `validate:"required,gte=1"`
	DecodedDepartmentId uint16  `validate:"required,gte=1"`
	Salary              float32 `validate:"required,gte=1"`
	DepartmentId        string  `validate:"required"`
	ManagerId           string  `validate:"required"`
	JobId               string  `validate:"required"`
	FirstName           string  `validate:"required,max=20"`
	LastName            string  `validate:"required,max=25"`
	EmailAddress        string  `validate:"required,max=100,email"`
	PhoneNumber         string  `validate:"required,max=20,numeric"`
	HireDate            string  `validate:"required"`
}
type updateEmployeeRequest struct {
	DecodedEmployeeID   uint16  `validate:"required,gte=1"`
	DecodedJobId        uint16  `validate:"required,gte=1"`
	DecodedManagerId    uint16  `validate:"required,gte=0"`
	DecodedDepartmentId uint16  `validate:"required,gte=1"`
	Salary              float32 `validate:"required,gte=1"`
	EmployeeID          string  `validate:"required"`
	DepartmentId        string  `validate:"required"`
	ManagerId           string  `validate:"required"`
	JobId               string  `validate:"required"`
	FirstName           string  `validate:"required,max=20"`
	LastName            string  `validate:"required,max=25"`
	EmailAddress        string  `validate:"required,max=100,email"`
	PhoneNumber         string  `validate:"required,max=20,numeric"`
	HireDate            string  `validate:"required"`
}
