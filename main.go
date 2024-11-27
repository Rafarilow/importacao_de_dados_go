package main

import (
	"fmt"
	"importador/db"
	"importador/utils"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	//Verifica se os parâmetros corretos foram passados
	if len(os.Args) < 6 {
		log.Fatal("Uso incorreto. Utilize: importador [arquivo xlsx] [host do db] [usuario do db] [senha do db] [nome do db]")
	}

	//Obtém os argumentos da linha de comando
	filePath := os.Args[1]
	dbHost := os.Args[2]
	dbUser := os.Args[3]
	dbPassword := os.Args[4]
	dbName := os.Args[5]

	//Conecta ao banco de dados usando as credenciais passadas
	err := db.ConnectDB(dbUser, dbPassword, dbHost, dbName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	//Lê os dados da planilha Excel
	contacts, err := readExcelFile(filePath)
	if err != nil {
		log.Fatalf("Erro ao ler a planilha: %v", err)
	}

	//Valida e insere cada contato no banco de dados
	for _, contact := range contacts {
		// Valida campos obrigatórios
		if err := utils.ValidateFields(contact); err != nil {
			fmt.Println("Erro na validação dos dados:", err)
			continue
		}

		//Valida o e-mail
		if !utils.ValidateEmail(contact["Email"]) {
			fmt.Println("E-mail inválido:", contact["Email"])
			continue
		}

		//Insere o contato no banco de dados
		if err := db.InsertContact(contact); err != nil {
			fmt.Println("Erro ao inserir no banco de dados:", err)
			continue
		}
		fmt.Println("Contato inserido:", contact["Nome"])
	}
}

// Função para ler a planilha
func readExcelFile(filePath string) ([]map[string]string, error) {
	//Abre o arquivo Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo Excel: %w", err)
	}
	defer f.Close()

	//Obtém as linhas da planilha (assumindo que a planilha tem apenas uma aba chamada "Sheet1")
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler as linhas da planilha: %w", err)
	}

	//Verifica se a planilha está vazia
	if len(rows) < 2 {
		return nil, fmt.Errorf("a planilha não contém dados suficientes")
	}

	//Determina a primeira linha como cabeçalhos
	headers := rows[0]
	var contacts []map[string]string

	//Processa as linhas a partir da segunda linha (dados)
	for _, row := range rows[1:] {
		if len(row) < len(headers) {
			// Linha de dados incompleta, pula
			continue
		}
		contact := make(map[string]string)
		for i, cell := range row {
			// Atribui os valores das células aos campos correspondentes (usando os cabeçalhos)
			contact[headers[i]] = cell
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}
