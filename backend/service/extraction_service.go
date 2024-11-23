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

func (s *ExtractionService) ExtractDataFromFile(ctx context.Context, file io.Reader, filename string) (*models.ExtractedDataCollection, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	fileContent, fileType, err := processFileByType(fileBytes, filename)

	fmt.Println("fileContent: ", fileType)

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

	// Updated schema to handle multiple invoices
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"invoices": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
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
				},
			},
			"missingFields": {Type: genai.TypeString},
		},
		Required: []string{"invoices"},
	}

	prompt := createPromptForFileType(fileType, fileContent)
	var resp *genai.GenerateContentResponse

	if fileType == "image" {
		fileExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
		prompt := createPromptForFileType("image", "")

		// Get the appropriate MIME type for the image
		mimeType := getImageMIMEType(filename)

		// Validate if the image type is supported
		if mimeType == "application/octet-stream" {
			return nil, fmt.Errorf("unsupported image type for file: %s", filename)
		}
		imgParts := []genai.Part{
			genai.ImageData(fileExt, fileBytes),
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

	var extractedDataCollection models.ExtractedDataCollection
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &extractedDataCollection); err != nil {
				return nil, fmt.Errorf("error parsing response: %v", err)
			}
			break
		}
	}

	if extractedDataCollection.MissingFields != "" {
		log.Printf("Missing fields detected in file %s: %s", filename, extractedDataCollection.MissingFields)
	}

	// Validate each extracted data entry
	// for i, data := range extractedDataCollection.Invoices {
	// 	if err := repository.ValidateExtractedData(&data); err != nil {
	// 		log.Printf("Validation error for invoice %d in file %s: %v", i+1, filename, err)
	// 		return nil, fmt.Errorf("validation failed for invoice %d: %v", i+1, err)
	// 	}
	// }

	// Save all extracted data to the database
	extractionRepo := repository.NewExtractionRepository()
	for i, data := range extractedDataCollection.Invoices {
		if err := extractionRepo.SaveExtractedData(ctx, &data, fmt.Sprintf("%s_invoice_%d", filename, i+1)); err != nil {
			return nil, fmt.Errorf("error saving extracted data for invoice %d: %v", i+1, err)
		}
	}

	return &extractedDataCollection, nil
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
	return fmt.Sprintf(`Extract all invoice, product, and customer information from this %s.
    For each invoice found, extract:
    - Invoice: serial number, date, and total amount
    - Products: name, quantity, unit price, tax and priceWithTax
    - Customer: name, phone number and totalPurchaseAmount
    
    If there are multiple invoices, include all of them in the response.
    
    IMPORTANT: If any fields are missing or cannot be extracted, please provide a detailed description 
    of which fields are missing in the "missingFields" field of the response. For example:
    - If customer phone number is missing: "Customer phone number not found in the document"
    - If product tax is missing: "Tax information missing for products"
    - If multiple fields are missing: "Missing fields: customer phone, product tax, invoice date"
    
    The response should be a JSON object containing:
    1. An "invoices" array with all extracted invoice data
    2. A "missingFields" string describing any missing information
    
    For partial data, still include what was found in the invoices array, and note the missing fields separately.
    
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

func getImageMIMEType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
