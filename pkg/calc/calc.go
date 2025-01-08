package calc

import (
	"math"
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
		return make([]string, 0), ErrEmptyExpression
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
				return make([]string, 0), ErrUnexpected(")")
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
			return make([]string, 0), ErrUnexpected("()")
		}
	}

	return rpm, nil
}

func Calc(expression string) (float64, error) {
	var stack []float64
	var a, b float64

	items, err := split(expression)
	if err != nil {
		return math.NaN(), err
	}

	rpm, err := toRPM(items)
	if err != nil {
		return math.NaN(), err
	}

	for _, item := range rpm {

		if strings.Contains(operations, item) {
			if len(stack) < 2 {
				return math.NaN(), ErrUnexpected(item)
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
					return math.NaN(), ErrDivisionByZero
				}
				stack = append(stack, a/b)
			default:
				return math.NaN(), ErrUnexpected(item)
			}
		} else {
			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return math.NaN(), ErrUnexpected(item)
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return math.NaN(), ErrInvalidExpression
	}

	return stack[0], nil
}
