package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"path/filepath"
	"strings"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/models"
	"github.com/TejasThombare20/backend/repository"
	"github.com/google/generative-ai-go/genai"
	"github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

type ExtractionService struct{}

func NewExtractionService() *ExtractionService {
	return &ExtractionService{}
}

// func (s *ExtractionService) ExtractDataFromFile(file io.Reader, fileType string) (map[string]interface{}, error) {
func (s *ExtractionService) ExtractDataFromFile(ctx context.Context, file io.Reader, filename string) (*models.ExtractedData, error) {

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	fileContent, fileType, err := processFileByType(fileBytes, filename)
	if err != nil {
		return nil, err
	}

	geminiClient, err := config.GetGeminiClient()
	if err != nil {
		return nil, fmt.Errorf("gemini client error: %v", err)
	}

	// Call Gemini API
	model := geminiClient.GenerativeModel("gemini-1.5-flash")

	model.ResponseMIMEType = "application/json"

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"invoice": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"serialNumber": {Type: genai.TypeString},
					"date":         {Type: genai.TypeString},
					"totalAmount":  {Type: genai.TypeNumber},
				},
				Required: []string{"serialNumber", "date", "totalAmount"},
			},
			"products": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name":         {Type: genai.TypeString},
						"quantity":     {Type: genai.TypeInteger},
						"unitPrice":    {Type: genai.TypeNumber},
						"tax":          {Type: genai.TypeNumber},
						"priceWithTax": {Type: genai.TypeNumber},
					},
					Required: []string{"name", "quantity", "unitPrice", "tax", "priceWithTax"},
				},
			},
			"customer": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"name":                {Type: genai.TypeString},
					"phoneNumber":         {Type: genai.TypeString},
					"totalPurchaseAmount": {Type: genai.TypeNumber},
				},
				Required: []string{"name", "phoneNumber", "totalPurchaseAmount"},
			},
		},
		Required: []string{"invoice", "products", "customer"},
	}

	prompt := createPromptForFileType(fileType, fileContent)

	var resp *genai.GenerateContentResponse
	if fileType == "image" {
		imgParts := []genai.Part{
			genai.ImageData("image/jpeg", []byte(fileContent)),
			genai.Text(prompt),
		}
		resp, err = model.GenerateContent(ctx, imgParts...)
	} else {
		resp, err = model.GenerateContent(ctx, genai.Text(prompt))
	}

	if err != nil {
		return nil, fmt.Errorf("error generating content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from model")
	}

	var extractedData models.ExtractedData
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &extractedData); err != nil {
				return nil, fmt.Errorf("error parsing response: %v", err)
			}
			break
		}
	}

	// Validate the extracted data
	if err := repository.ValidateExtractedData(&extractedData); err != nil {
		// Log the validation error
		log.Printf("Validation error for file %s: %v", filename, err)

		// You can choose to return an error or handle it differently based on your requirements
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// Save the extracted data to the database
	extractionRepo := repository.NewExtractionRepository()
	if err := extractionRepo.SaveExtractedData(ctx, &extractedData, filename); err != nil {
		return nil, fmt.Errorf("error saving extracted data: %v", err)
	}

	return &extractedData, nil

}

func processFileByType(fileBytes []byte, filename string) (string, string, error) {
	fileExt := strings.ToLower(filepath.Ext(filename))

	switch fileExt {
	case ".pdf":
		content, err := extractPDFText(fileBytes)
		return content, "pdf", err
	case ".xlsx", ".xls":
		content, err := extractExcelText(fileBytes)
		return content, "excel", err
	case ".jpg", ".jpeg", ".png":
		content := base64.StdEncoding.EncodeToString(fileBytes)
		return content, "image", nil
	default:
		return "", "", fmt.Errorf("unsupported file type: %s", fileExt)
	}
}

func extractExcelText(fileBytes []byte) (string, error) {
	// Create a new buffer from bytes
	reader := bytes.NewReader(fileBytes)

	// Open Excel file
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content strings.Builder

	// Get all sheet names
	sheets := f.GetSheetList()

	// Iterate through each sheet
	for _, sheet := range sheets {
		// Get all cells from the sheet
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}

		content.WriteString(fmt.Sprintf("Sheet: %s\n", sheet))

		// Convert rows to text
		for _, row := range rows {
			content.WriteString(strings.Join(row, "|") + "\n")
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}

func createPromptForFileType(fileType string, content string) string {
	return fmt.Sprintf(`Extract invoice, product, and customer information from this %s.
    Focus on finding:
    - Invoice: serial number, date, and total amount
    - Products: name, quantity, unit price, tax and priceWithTax
    - Customer: name, phone number and totalPurchaseAmount
    
    Respond with a JSON object matching the specified schema.
    
    Content: %s`, fileType, content)
}

func extractPDFText(fileBytes []byte) (string, error) {
	// Create a temporary reader from bytes
	reader := bytes.NewReader(fileBytes)

	// Parse PDF
	pdfReader, err := pdf.NewReader(reader, int64(len(fileBytes)))
	if err != nil {
		return "", err
	}

	var content strings.Builder
	for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
		page := pdfReader.Page(pageNum)
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		content.WriteString(text)
	}

	return content.String(), nil
}
