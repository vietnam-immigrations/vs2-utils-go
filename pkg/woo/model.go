package woo

import "fmt"

type Billing struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Company           string `json:"company"`
	Address1          string `json:"address_1"`
	Address2          string `json:"address_2"`
	City              string `json:"city"`
	State             string `json:"state"`
	Postcode          string `json:"postcode"`
	Country           string `json:"country"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Vs2SecondaryEmail string `json:"vs2_secondary_email"`
}

type Shipping struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
}

type Meta struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Value        string `json:"value"`
	DisplayKey   string `json:"display_key"`
	DisplayValue string `json:"display_value"`
}

type Metas []Meta

func (ms Metas) GetMeta(key string) (*Meta, error) {
	for _, meta := range ms {
		if meta.Key == key {
			return &meta, nil
		}
	}
	return nil, fmt.Errorf("meta data [%s] not found", key)
}

type Image struct {
	ID  string `json:"id"`
	Src string `json:"src"`
}

type LineItem struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	ProductID   int           `json:"product_id"`
	VariationID int           `json:"variation_id"`
	Quantity    int           `json:"quantity"`
	TaxClass    string        `json:"tax_class"`
	Subtotal    string        `json:"subtotal"`
	SubtotalTax string        `json:"subtotal_tax"`
	Total       string        `json:"total"`
	TotalTax    string        `json:"total_tax"`
	Taxes       []interface{} `json:"taxes"`
	MetaData    Metas         `json:"meta_data"`
	SKU         string        `json:"sku"`
	Price       int           `json:"price"`
	Image       Image         `json:"image"`
	ParentName  interface{}   `json:"parent_name"`
}

type FeeLine struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	TaxClass  string        `json:"tax_class"`
	TaxStatus string        `json:"tax_status"`
	Amount    string        `json:"amount"`
	Total     string        `json:"total"`
	TotalTax  string        `json:"total_tax"`
	Taxes     []interface{} `json:"taxes"`
	MetaData  Metas         `json:"meta_data"`
}

type Vs2CheckoutCustomFields struct {
	Vs2ArrivalDate      string `json:"vs2_arrival_date"`
	Vs2Checkpoint       string `json:"vs2_checkpoint"`
	Vs2ProcessingTime   string `json:"vs2_processing_time"`
	Vs2FastTrack        string `json:"vs2_fast_track"`
	Vs2CarPickup        string `json:"vs2_car_pickup"`
	Vs2Flight           string `json:"vs2_flight"`
	Vs2CarPickupAddress string `json:"vs2_car_pickup_address"`
}

type Order struct {
	ID                 int                     `json:"id"`
	ParentID           int                     `json:"parent_id"`
	Status             string                  `json:"status"`
	Currency           string                  `json:"currency"`
	Version            string                  `json:"version"`
	PricesIncludeTax   bool                    `json:"prices_include_tax"`
	DateCreated        string                  `json:"date_created"`
	DateModified       string                  `json:"date_modified"`
	DiscountTotal      string                  `json:"discount_total"`
	DiscountTax        string                  `json:"discount_tax"`
	ShippingTotal      string                  `json:"shipping_total"`
	ShippingTax        string                  `json:"shipping_tax"`
	CartTax            string                  `json:"cart_tax"`
	Total              string                  `json:"total"`
	TotalTax           string                  `json:"total_tax"`
	CustomerID         int                     `json:"customer_id"`
	OrderKey           string                  `json:"order_key"`
	Billing            Billing                 `json:"billing"`
	Shipping           Shipping                `json:"shipping"`
	PaymentMethod      string                  `json:"payment_method"`
	PaymentMethodTitle string                  `json:"payment_method_title"`
	TransactionId      string                  `json:"transaction_id"`
	CustomerIpAddress  string                  `json:"customer_ip_address"`
	CustomerUserAgent  string                  `json:"customer_user_agent"`
	CreatedVia         string                  `json:"created_via"`
	CustomerNote       string                  `json:"customer_note"`
	DateCompleted      string                  `json:"date_completed"`
	DatePaid           string                  `json:"date_paid"`
	CartHash           string                  `json:"cart_hash"`
	Number             string                  `json:"number"`
	MetaData           Metas                   `json:"meta_data"`
	LineItems          []LineItem              `json:"line_items"`
	TaxLines           []interface{}           `json:"tax_lines"`
	ShippingLines      []interface{}           `json:"shipping_lines"`
	FeeLines           []FeeLine               `json:"fee_lines"`
	CouponLines        []interface{}           `json:"coupon_lines"`
	Refunds            []interface{}           `json:"refunds"`
	PaymentUrl         string                  `json:"payment_url"`
	IsEditable         bool                    `json:"is_editable"`
	NeedsPayment       bool                    `json:"needs_payment"`
	NeedsProcessing    bool                    `json:"needs_processing"`
	DateCreatedGmt     string                  `json:"date_created_gmt"`
	DateModifiedGmt    string                  `json:"date_modified_gmt"`
	DateCompletedGmt   interface{}             `json:"date_completed_gmt"`
	DatePaidGmt        string                  `json:"date_paid_gmt"`
	Vs2EvisaCheckout   Vs2CheckoutCustomFields `json:"vs2_evisa_checkout"`
	CurrencySymbol     string                  `json:"currency_symbol"`
}
