package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
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

type Job struct {
	Employee   Employee
	ResultCalc float64
}

func worker(id int, jobs <-chan Job, results chan<- Job) {
	for job := range jobs {
		fmt.Println("worker", id, "stard job", job.Employee.ID)

		// convert salary and age to float for calculation
		salary, err1 := job.Employee.EmployeeSalary.Float64()
		age, err2 := job.Employee.EmployeeAge.Float64()

		if err1 != nil {
			fmt.Printf("Worker %d error parsing salary for employee %s: %v\n", id, job.Employee.EmployeeName, err1)
			job.ResultCalc = 0
		} else if err2 != nil {
			fmt.Printf("Worker %d error parsing age for employee %s: %v\n", id, job.Employee.EmployeeName, err2)
			job.ResultCalc = 0
		} else if age != 0 {
			job.ResultCalc = salary / age
		} else {
			job.ResultCalc = 0 // handle division by zero
		}

		fmt.Printf("Worker %d processed job for employee %s with result %f\n", id, job.Employee.EmployeeName, job.ResultCalc)
		results <- job
		fmt.Println("worker", id, "finished job", job.Employee.ID)
	}
}

func fetchEmployeeData() ([]Employee, error) {
	url := "https://dummy.restapiexample.com/api/v1/employees"

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return nil, err
	}

	// Define a variable to hold the full response
	var response Response

	// Unmarshal the JSON data into the Response struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return nil, err
	}

	return response.Data, nil

}

func main() {
	// Fetch data from API
	employees, err := fetchEmployeeData()
	if err != nil {
		log.Fatalf("Failed to fetch employees: %s", err)
	}

	const numWorkers = 5

	jobs := make(chan Job, len(employees))
	results := make(chan Job, len(employees))
	var wg sync.WaitGroup

	// Start workers pool
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	// Push jobs to worker pool
	for _, emp := range employees {
		jobs <- Job{Employee: emp}
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	// Collect results
	for result := range results {
		fmt.Printf("Employee %s (Salary: %s, Age: %s) -> Salary/Age: %f\n",
			result.Employee.EmployeeName, result.Employee.EmployeeSalary, result.Employee.EmployeeAge, result.ResultCalc)
	}
}
