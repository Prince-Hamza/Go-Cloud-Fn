package Structures

type ITScopeInfo struct {
	Products []struct {
		Id                    *string
		Sku                   string  `json:"Sku"`
		Category              *string `json:"category"`
		Price                 *string `json:"price"`
		Stock                 *string `json:"stock"`
		StockStatus           *string `json:"stockStatus"`
		StandardHtmlDatasheet *string `json:"standardHtmlDatasheet"`
		Extra                 *string `json:"extra"`
		Attributes            []struct {
			ID      *int     `json:"id"`
			Name    string   `json:"name"`
			Options []string `json:"options"`
		} `json:"attributes"`

		Categories []struct {
			ID int `json:"id"`
		} `json:"categories"`
		Images []struct {
			Src string `json:"src"`
		} `json:"images"`
	} `json:"products"`
}

type WoocommerceInfo struct {
	FieldsInResponse [3]string `json:"fields_in_response"`
	Type             string    `json:"type"`
	RegularPrice     string    `json:"regular_price"`
	Price            string    `json:"price"`
	ManageStock      bool      `json:"manage_stock"`
	StockStatus      string    `json:"stock_status"`
	StockQuantity    string    `json:"stock_quantity"`
	Extra            *string   `json:"extra"`
	Attributes       []struct {
		ID      *int     `json:"id"`
		Name    string   `json:"name"`
		Options []string `json:"options"`
	} `json:"attributes"`
	Categories []struct {
		ID int `json:"id"`
	} `json:"categories"`
	Images []struct {
		Src string `json:"src"`
	} `json:"images"`
}

type FieldsInResponse struct {
}
