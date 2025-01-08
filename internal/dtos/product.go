package dtos

type CreateProductDto struct {
	Name           string  `json:"name" example:"Shampoo Anticaspa"`
	Brand          string  `json:"brand" example:"Head & Shoulders"`
	Unit           string  `json:"unit" example:"ml"`
	PackageCount   float64 `json:"package_count" example:"32"`
	UnitPerPackage float64 `json:"unit_per_package" example:"500"`
	LowStockAlert  float64 `json:"low_stock_alert" example:"100"`
	UnityPrice     float64 `json:"unity_price" example:"10000"`
}

type UpdateProductDto struct {
	Name          string  `json:"name" example:"Shampoo Anticaspa"`
	Brand         string  `json:"brand" example:"Head & Shoulders"`
	Unit          string  `json:"unit" example:"ml"`
	LowStockAlert float64 `json:"low_stock_alert" example:"100"`
}

type RestockProductDto struct {
	PackageCount   float64 `json:"package_count" example:"32"`
	UnitPerPackage float64 `json:"unit_per_package" example:"500"`
	Reason         string  `json:"reason" example:"Reinventario"` // Razón del movimiento (opcional)
	UnityPrice     float64 `json:"unity_price" example:"10000"`   // Precio unitario
}

type GetProductDto struct {
	ID       uint    `json:"id" example:"1"`
	Name     string  `json:"name" example:"Shampoo Anticaspa"`
	Brand    string  `json:"brand" example:"Head & Shoulders"`
	Unit     string  `json:"unit" example:"ml"`
	Quantity float64 `json:"quantity" example:"400"`
}
