package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SerialNumber  *string            `bson:"serialNumber,omitempty" json:"serialNumber,omitempty"`
	TotalAmount   *float64           `bson:"totalAmount,omitempty" json:"totalAmount,omitempty"`
	Date          *time.Time         `bson:"date,omitempty" json:"date,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	ExtractedFrom *string            `bson:"extracted_from,omitempty" json:"extracted_from,omitempty"`
}

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Quantity     *int               `bson:"quantity,omitempty" json:"quantity,omitempty"`
	UnitPrice    *float64           `bson:"unitPrice,omitempty" json:"unitPrice,omitempty"`
	Tax          *float64           `bson:"tax,omitempty" json:"tax,omitempty"`
	PriceWithTax *float64           `bson:"priceWithTax,omitempty" json:"priceWithTax,omitempty"`
	InvoiceID    primitive.ObjectID `bson:"invoice_id,omitempty" json:"invoice_id,omitempty"`
}

type Customer struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name                *string            `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber         *string            `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
	TotalPurchaseAmount *float64           `bson:"totalPurchaseAmount,omitempty" json:"totalPurchaseAmount,omitempty"`
	InvoiceID           primitive.ObjectID `bson:"invoice_id,omitempty" json:"invoice_id,omitempty"`
}

func (p *Product) ToUpdateMapProduct() map[string]interface{} {
	update := make(map[string]interface{})

	if p.Name != "" {
		update["name"] = p.Name
	}
	if p.Quantity != nil {
		update["quantity"] = *p.Quantity
	}
	if p.UnitPrice != nil {
		update["unitPrice"] = *p.UnitPrice
	}
	if p.Tax != nil {
		update["tax"] = *p.Tax
	}
	if p.PriceWithTax != nil {
		update["priceWithTax"] = *p.PriceWithTax
	}
	if !p.InvoiceID.IsZero() {
		update["invoice_id"] = p.InvoiceID
	}

	return update
}

func IntPtr(v int) *int {
	return &v
}

// Helper function to create a pointer to a float64
func Float64Ptr(v float64) *float64 {
	return &v
}

func (c *Customer) ToUpdateMapCustomer() map[string]interface{} {
	updates := make(map[string]interface{})

	if c.Name != nil {
		updates["name"] = *c.Name
	}
	if c.PhoneNumber != nil {
		updates["phoneNumber"] = *c.PhoneNumber
	}
	if c.TotalPurchaseAmount != nil {
		updates["totalPurchaseAmount"] = *c.TotalPurchaseAmount
	}
	if !c.InvoiceID.IsZero() {
		updates["invoice_id"] = c.InvoiceID
	}

	return updates
}

func (i *Invoice) ToUpdateMapInvoice() map[string]interface{} {
	updates := make(map[string]interface{})

	if i.SerialNumber != nil {
		updates["serialNumber"] = *i.SerialNumber
	}
	if i.TotalAmount != nil {
		updates["totalAmount"] = *i.TotalAmount
	}
	if i.Date != nil {
		updates["date"] = *i.Date
	}

	if i.ExtractedFrom != nil {
		updates["extracted_from"] = *i.ExtractedFrom
	}

	return updates
}
