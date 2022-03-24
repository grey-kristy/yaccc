package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"kristy/convertor"
)

func main() {
	amount, from, to := get_args()
	price := convertor.Convert(amount, from, to)
	if price != 0 {
		fmt.Printf("%f %s is %f %s\n", amount, from, price, to)
	}
}

func print_usage() {
	fmt.Println("Yet Another Crypto Currency Converter\n")
	fmt.Println("Usage: ")
	fmt.Println("yaccc <AMOUNT> <FROM_CURRENCY_SYMBOL> <TO_CURRENCY_SYMBOL>\n")
	fmt.Println("For example:")
	fmt.Println("yaccc.py 0.5 BTC USD\n")
}

func get_args() (float64, string, string) {
	if len(os.Args) < 4 {
		print_usage()
		os.Exit(4)
	}
	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Amount should be integer of float, not \"%s\"\n", os.Args[1])
		os.Exit(4)
	}
	return amount, strings.ToUpper(os.Args[2]), strings.ToUpper(os.Args[3])
}
