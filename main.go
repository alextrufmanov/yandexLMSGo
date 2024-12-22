package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const operations string = "+-*/"

func priority(operation string) int {
	switch operation {
	case "+":
		return 1
	case "-":
		return 1
	case "*":
		return 2
	case "/":
		return 2
	default:
		return -1
	}
}

func split(expression string) ([]string, error) {

	expression = strings.ReplaceAll(expression, " ", "")

	if len(expression) == 0 {
		return make([]string, 0), errors.New("invalid input, expression is empty")
	}

	if string(expression[0]) == "-" {
		expression = "0" + expression
	}
	expression = strings.ReplaceAll(expression, "(-", "(0-")

	for _, r := range operations + "()" {
		expression = strings.ReplaceAll(expression, string(r), " "+string(r)+" ")
	}

	return strings.Fields(expression), nil
}

func toRPM(items []string) ([]string, error) {
	var rpm []string
	var stack []string

	for _, item := range items {
		switch {
		case strings.Contains(operations, item):
			for ; len(stack) > 0; stack = stack[:len(stack)-1] {
				if priority(stack[len(stack)-1]) >= priority(item) && stack[len(stack)-1] != "(" {
					rpm = append(rpm, stack[len(stack)-1])
				} else {
					break
				}
			}
			stack = append(stack, item)
		case item == "(":
			stack = append(stack, item)
		case item == ")":
			for ; len(stack) > 0; stack = stack[:len(stack)-1] {
				if stack[len(stack)-1] != "(" {
					rpm = append(rpm, stack[len(stack)-1])
				} else {
					break
				}
			}
			if len(stack) == 0 {
				return make([]string, 0), errors.New("invalid input, unexpected ')'")
			}
			stack = stack[:len(stack)-1]
		default:
			rpm = append(rpm, item)
		}
	}

	for ; len(stack) > 0; stack = stack[:len(stack)-1] {
		if stack[len(stack)-1] != "(" {
			rpm = append(rpm, stack[len(stack)-1])
		} else {
			return make([]string, 0), errors.New("invalid input, unexpected '('")
		}
	}

	return rpm, nil
}

func Calc(expression string) (float64, error) {
	var stack []float64
	var a, b float64

	items, err := split(expression)
	if err != nil {
		return math.NaN(), errors.New(err.Error())
	}

	rpm, err := toRPM(items)
	if err != nil {
		return math.NaN(), errors.New(err.Error())
	}

	for _, item := range rpm {

		if strings.Contains(operations, item) {
			if len(stack) < 2 {
				return math.NaN(), errors.New(fmt.Sprintf("invalid input, unexpected '%s'", item))
			}
			a = stack[len(stack)-2]
			b = stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			switch {
			case item == "+":
				stack = append(stack, a+b)
			case item == "-":
				stack = append(stack, a-b)
			case item == "*":
				stack = append(stack, a*b)
			case item == "/":
				if b == 0 {
					return math.NaN(), errors.New("invalid input, division by zero")
				}
				stack = append(stack, a/b)
			default:
				return math.NaN(), errors.New(fmt.Sprintf("invalid input, undefined operation '%s'", item))
			}
		} else {
			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return math.NaN(), errors.New(fmt.Sprintf("invalid input, invalid value '%s'", item))
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return math.NaN(), errors.New(fmt.Sprintf("invalid input, expression"))
	}

	return stack[0], nil
}

type RequestBody struct {
	Expression string `json:"expression"`
}

type AnswerBody struct {
	Result string `json:"result"`
}

func ErrorNotFound(w http.ResponseWriter) {
	http.Error(w, "404 not found.", http.StatusNotFound)
}

func CalcError(w http.ResponseWriter) {
	http.Error(w, "{\n    \"error\": \"Expression is not valid\"\n}", http.StatusUnprocessableEntity)
}

func InternalError(w http.ResponseWriter) {
	http.Error(w, "{\n    \"error\": \"Internal server error\"\n}", http.StatusInternalServerError)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody RequestBody
	var answerBody AnswerBody

	if strings.ToUpper(r.Method) != "POST" {
		ErrorNotFound(w)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil {
		if json.Unmarshal(bodyBytes, &requestBody) == nil {
			calcResult, err := Calc(requestBody.Expression)
			if err == nil {
				answerBody.Result = fmt.Sprint(calcResult)
				outData, _ := json.Marshal(answerBody)
				fmt.Fprintf(w, string(outData))
			} else {
				CalcError(w)
			}
			return
		}
	}
	InternalError(w)
}

func StartServer() {
	fmt.Println("Listen localhost:80")
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	http.ListenAndServe(":80", nil)
}

func main() {
	StartServer()
}
