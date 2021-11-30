package Structures

type ITScopeInfo struct {
	Products []struct {
		ManufacturerSKU   string `json:"manufacturerSKU"`
		ProductSubType    string `json:"productSubType"`
		SupplierPriceInfo struct {
			Price string `json:"price"`
		} `json:"supplierPriceInfo"`
		ProductStockInfo struct {
			Stock           string `json:"stock"`
			StockStatusText string `json:"stockStatusText"`
		} `json:"productStockInfo"`
		AggregatedStock           string `json:"aggregatedStock"`
		AggregatedStockStatusText string `json:"aggregatedStockStatusText"`
		Extra                     string `json:"extra"`
	} `json:"products"`
}

type WoocommerceInfo struct {
	FieldsInResponse []string `json:"fields_in_response"`
	Type             string   `json:"type"`
	RegularPrice     string   `json:"regular_price"`
	Price            string   `json:"price"`
	ManageStock      string   `json:"manage_stock"`
	StockStatus      string   `json:"stock_status"`
	StockQuantity    string   `json:"stock_quantity"`
	Attributes       []struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Options []string `json:"options"`
	} `json:"attributes"`
}
