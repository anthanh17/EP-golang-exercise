package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Employee struct {
	ID             json.Number `json:"id"`
	EmployeeName   string      `json:"employee_name"`
	EmployeeSalary json.Number `json:"employee_salary"`
	EmployeeAge    json.Number `json:"employee_age"`
	ProfileImage   string      `json:"profile_image"`
}

type Response struct {
	Status  string     `json:"status"`
	Data    []Employee `json:"data"`
	Message string     `json:"message"`
}

func main() {
	url := "https://dummy.restapiexample.com/api/v1/employees"

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// fmt.Println("Body - String:", string(body))

	// Define a variable to hold the full response
	var response Response

	// Unmarshal the JSON data into the Response struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	// Print the employees
	for _, employee := range response.Data {
		fmt.Printf("{ Id: %s, Name: %s, Salary: %s, Age: %s, Profile Image: %s }\n",
			employee.ID, employee.EmployeeName, employee.EmployeeSalary, employee.EmployeeAge, employee.ProfileImage)
	}
}
