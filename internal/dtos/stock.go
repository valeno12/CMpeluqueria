package dtos

type StockMovementDto struct {
	ID             uint     `json:"id" example:"1"`
	ProductID      uint     `json:"product_id" example:"10"`
	ProductName    string   `json:"product_name" example:"Shampoo"`
	ProductBrand   string   `json:"product_brand" example:"Pantene"`
	ProductUnit    string   `json:"product_unit" example:"lt"`
	Quantity       float64  `json:"quantity" example:"20.5"`
	PackageCount   *float64 `json:"package_count,omitempty" example:"5"`    // Mostrar solo si es entrada
	UnitPerPackage *float64 `json:"unit_per_package,omitempty" example:"4"` // Mostrar solo si es entrada
	UnityPrice     *float64 `json:"unity_price,omitempty" example:"15000"`  // Mostrar solo si es entrada
	Reason         string   `json:"reason" example:"Compra de stock"`
	CreatedAt      string   `json:"created_at" example:"30/09/2025 15:30"`
}
