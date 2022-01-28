package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ResponseBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Operation struct {
	Operand1  float64 `json:"operand1"`
	Operand2  float64 `json:"operand2"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

const DefaultPort int = 8080

var history []Operation

func createJsonResponseWithData(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ResponseBody{message, data})
}

func createSimpleJsonResponse(w http.ResponseWriter, message string, statusCode int) {
	createJsonResponseWithData(w, message, nil, statusCode)
}

func executeCalculator(operand1, operand2 float64, operation string) (result float64, err error) {
	switch operation {
	case "sum":
		result = operand1 + operand2
	case "sub":
		result = operand1 - operand2
	case "mul":
		result = operand1 * operand2
	case "div":
		if operand2 == 0.0 {
			err = errors.New("Cannot divide by zero")
			return
		}
		result = operand1 / operand2
	default:
		err = errors.New(operation + " is not a valid operation")
		return
	}
	return
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		createSimpleJsonResponse(w, "Missing path parameters in /calc/{operation}/{operand1}/{operand2}", http.StatusBadRequest)
		return
	}

	operation := pathParts[2]
	operand1, conversionErr1 := strconv.ParseFloat(pathParts[3], 64)
	if conversionErr1 != nil {
		createSimpleJsonResponse(w, conversionErr1.Error(), http.StatusBadRequest)
		return
	}
	operand2, conversionErr2 := strconv.ParseFloat(pathParts[4], 64)
	if conversionErr2 != nil {
		createSimpleJsonResponse(w, conversionErr2.Error(), http.StatusBadRequest)
		return
	}

	result, calculatorError := executeCalculator(operand1, operand2, operation)
	if calculatorError != nil {
		createSimpleJsonResponse(w, calculatorError.Error(), http.StatusBadRequest)
		return
	}

	operationInfo := Operation{operand1, operand2, operation, result}
	history = append(history, operationInfo)

	createJsonResponseWithData(w, "Operation successful", operationInfo, 200)
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	if len(history) == 0 {
		createSimpleJsonResponse(w, "Calculator history is empty", 200)
		return
	}
	createJsonResponseWithData(w, "Found history from oldest to most recent", history, 200)
}

func getServerPort() (port int) {
	value := os.Getenv("SERVER_PORT")
	if value == "" {
		return DefaultPort
	}
	portValue, err := strconv.Atoi(value)
	if err != nil {
		return DefaultPort
	}
	return portValue
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/calc/", calcHandler)
	http.HandleFunc("/calc/history", historyHandler)
	port := getServerPort()
	log.Println("Server started running at port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}