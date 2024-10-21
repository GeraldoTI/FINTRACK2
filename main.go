package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate-pdf", generatePDFHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ec2Cost := r.FormValue("ec2Cost")
	rdsCost := r.FormValue("rdsCost")
	vpcCost := r.FormValue("vpcCost")
	availableBudget := r.FormValue("availableBudget")

	totalCost := calculateTotal(ec2Cost, rdsCost, vpcCost)

	availableBudgetFloat, err := strconv.ParseFloat(availableBudget, 64)
	if err != nil {
		http.Error(w, "Invalid available budget", http.StatusBadRequest)
		return
	}

	totalCostFloat, err := strconv.ParseFloat(strings.Trim(totalCost, "$"), 64)
	if err != nil {
		http.Error(w, "Invalid total cost", http.StatusBadRequest)
		return
	}

	remainingBudget := availableBudgetFloat - totalCostFloat

	// Gerar conteúdo do PDF como texto simples
	report := generateReport(ec2Cost, rdsCost, vpcCost, totalCost, availableBudgetFloat, remainingBudget)

	// Salvar o relatório em um arquivo de texto (pode ser renomeado para .pdf)
	err = os.WriteFile("output.txt", []byte(report), 0644)
	if err != nil {
		http.Error(w, "Error creating report file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve o arquivo de relatório para download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
	http.ServeFile(w, r, "output.txt")

	// Limpar o arquivo após o uso
	os.Remove("output.txt")
}

func calculateTotal(ec2, rds, vpc string) string {
	ec2Cost, _ := strconv.ParseFloat(ec2, 64)
	rdsCost, _ := strconv.ParseFloat(rds, 64)
	vpcCost, _ := strconv.ParseFloat(vpc, 64)
	total := ec2Cost + rdsCost + vpcCost
	return fmt.Sprintf("$%.2f", total)
}

func generateReport(ec2, rds, vpc string, totalCost string, availableBudget float64, remainingBudget float64) string {
	return fmt.Sprintf(
		"Cloud Cost Report\n\n"+
			"EC2 Provisioned Cost: $%s\n"+
			"RDS Provisioned Cost: $%s\n"+
			"VPC Provisioned Cost: $%s\n\n"+
			"Total Provisioned Cost: %s\n"+
			"Available Budget: $%.2f\n"+
			"Remaining Budget: $%.2f\n\n"+
			"Graphical representation (text-based):\n\n"+
			"EC2: %s\n"+
			"RDS: %s\n"+
			"VPC: %s\n"+
			"Available Budget: %s\n",
		ec2, rds, vpc, totalCost,
		availableBudget, remainingBudget,
		generateBarGraph(ec2),
		generateBarGraph(rds),
		generateBarGraph(vpc),
		generateBarGraph(fmt.Sprintf("%.2f", availableBudget)),
	)
}

func generateBarGraph(valueStr string) string {
	maxVal := 100.0 // Valor máximo para normalização
	val, _ := strconv.ParseFloat(valueStr, 64)
	barLength := int((val / maxVal) * 40) // 40 caracteres de largura
	return fmt.Sprintf("%s: %s", valueStr, strings.Repeat("█", barLength))
}
