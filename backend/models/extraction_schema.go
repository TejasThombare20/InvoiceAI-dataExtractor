package models

type ExtractedData struct {
	Invoice  InvoiceData   `json:"invoice"`
	Products []ProductData `json:"products"`
	Customer CustomerData  `json:"customer"`
}

type ExtractedDataCollection struct {
	Invoices      []ExtractedData `json:"invoices"`
	MissingFields string          `json:"missingFields,omitempty"`
}

// type ExtractedData struct {
// 	Invoice  InvoiceData   `json:"invoice"`
// 	Products []ProductData `json:"products"`
// 	Customer CustomerData  `json:"customer"`
// }

type InvoiceData struct {
	SerialNumber string  `json:"serialNumber"`
	Date         string  `json:"date"`
	TotalAmount  float64 `json:"totalAmount"`
}

type ProductData struct {
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unitPrice"`
	Tax          float64 `json:"tax"`
	PriceWithTax float64 `json:"priceWithTax"`
}

type CustomerData struct {
	Name                string  `json:"name"`
	PhoneNumber         string  `json:"phoneNumber"`
	TotalPurchaseAmount float64 `json:"totalPurchaseAmount"`
}
