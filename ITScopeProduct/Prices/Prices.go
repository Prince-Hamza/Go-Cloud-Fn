package Price

import (
	"fmt"
)

type Prices struct{}

func (prices Prices) GetFinalPrices() {
	fmt.Println("prices package")
}
