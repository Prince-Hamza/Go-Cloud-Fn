package Product

import (
	//"encoding/json"
	//"fmt"
)

type Product struct {
	price              string
	regular_price      string
	manage_stock       bool
	stock_quantity     int
	stock_status       string
	fields_in_response [3]string
}

type Attributes struct {
	name    string
	id      string
	options AttributeOptions
}

type AttributeOptions struct {
	options [0]string
}

