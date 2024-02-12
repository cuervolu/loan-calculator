package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

var (
	typeOfCalc     string
	loanPrincipal  int
	periods        int
	monthlyPayment float64
	interest       float64
)
var allowedTypes = map[string]bool{
	"annuity": true,
	"diff":    true,
}

const errMessage = "Incorrect parameters"

func main() {
	setup()
	typeOfCalculation()
}

func typeOfCalculation() {
	switch typeOfCalc {
	case "annuity":
		if loanPrincipal == 0 {
			loanPrincipal = int(calculatePrincipal(monthlyPayment, interest, float64(periods)))
			fmt.Println("Your loan principal = ", loanPrincipal, "!")
		}
		if monthlyPayment == 0 {
			monthlyPayment = calculateMonthlyPayment(float64(loanPrincipal), interest, float64(periods))
			fmt.Println("Your monthly payment = ", monthlyPayment, "!")
		}
		if periods == 0 {
			periods = int(math.Ceil(math.Log(monthlyPayment/(monthlyPayment-(interest/100/12)*float64(loanPrincipal))) / math.Log(1+(interest/100/12))))
			fmt.Println("It will take ", periods/12, " months to repay this loan!")
		}
	case "diff":
		calculateDiffPayment()
	default:
		fmt.Println(errMessage)
		os.Exit(1)
	}
	var overpayment = calculateOverpayment(monthlyPayment, float64(loanPrincipal), float64(periods))
	if overpayment < 0 {
		fmt.Println(errMessage)
		os.Exit(1)
	}
	fmt.Println("\nOverpayment = ", overpayment)
}

func calculateDiffPayment() {
	totalPayment := 0.0
	for i := 1; i <= periods; i++ {
		payment := math.Ceil(float64(loanPrincipal)/float64(periods) + (interest/100/12)*(float64(loanPrincipal)-float64(loanPrincipal)*(float64(i)-1)/float64(periods)))
		totalPayment += payment
		fmt.Println("Month ", i, ": payment is ", payment)
	}
	fmt.Println("\nOverpayment = ", int(totalPayment)-loanPrincipal)
}

func validateInputs() {
	if !allowedTypes[typeOfCalc] {
		fmt.Println(errMessage)
		os.Exit(1)
	}
	if loanPrincipal < 0 || periods < 0 || monthlyPayment < 0 || interest < 0 {
		fmt.Println(errMessage)
		os.Exit(1)
	}

	if loanPrincipal == 0 && periods == 0 && monthlyPayment == 0 {
		fmt.Println(errMessage)
		os.Exit(1)
	}

	if loanPrincipal == 0 && typeOfCalc == "diff" {
		fmt.Println(errMessage)
		os.Exit(1)
	}
	//
	//if monthlyPayment == 0 && typeOfCalc == "diff" {
	//	fmt.Println(errMessage)
	//	os.Exit(1)
	//}

	if flag.NFlag() != 4 {
		fmt.Println(errMessage)
		os.Exit(1)
	}
}

func calculatePrincipal(payment, interest, numPayments float64) float64 {
	return payment / ((interest / 100 / 12) * math.Pow(1+(interest/100/12), numPayments) / (math.Pow(1+(interest/100/12), numPayments) - 1))
}

func calculateMonthlyPayment(principal, interest, numPayments float64) float64 {
	payment := principal * (interest / 100 / 12) * math.Pow(1+(interest/100/12), numPayments) /
		(math.Pow(1+(interest/100/12), numPayments) - 1)

	return math.Ceil(payment) // Round before returning
}

func calculateOverpayment(payment, principal, periods float64) int32 {

	totalPaid := payment * periods

	overpayment := totalPaid - principal

	return int32(math.Ceil(overpayment))
}

func setup() {
	flag.StringVar(&typeOfCalc, "type", "", "Type of calculation: annuity or diff")
	flag.IntVar(&loanPrincipal, "principal", 0, "Loan principal")
	flag.IntVar(&periods, "periods", 0, "Number of months needed to repay the loan")
	flag.Float64Var(&monthlyPayment, "payment", 0, "Monthly payment amount")
	flag.Float64Var(&interest, "interest", 0, "Annual interest rate")
	flag.Parse()
	validateInputs()
}
