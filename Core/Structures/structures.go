package Structures

type ITScopeInfo struct {
	Products []struct {
		Id          *string `json:"id"`
		Title       *string `json:"title"`
		Description *string `json:"description"`

		Sku                   string  `json:"sku"`
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
		Images []string `json:"images"`
		WpImages *[]struct {
           Src string
		} `json:"wpImages"`
	} `json:"products"`
}

type WoocommerceInfo struct {
	FieldsInResponse [3]string `json:"fields_in_response"`
	Name             *string   `json:"name"`
	Description      *string   `json:"description"`
	Sku              string    `json:"sku"`
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
	// Images []string `json:"images"`
	Images *[]struct {
		Src string
	 } `json:"images"`
}

type FieldsInResponse struct {
}
