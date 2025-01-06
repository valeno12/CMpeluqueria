package dtos

type CreateProductDto struct {
	Name           string  `json:"name"`
	Brand          string  `json:"brand"`
	Unit           string  `json:"unit"`
	PackageCount   float64 `json:"package_count"`
	UnitPerPackage float64 `json:"unit_per_package"`
	LowStockAlert  float64 `json:"low_stock_alert"`
	UnityPrice     float64 `json:"unity_price"`
}

type UpdateProductDto struct {
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Unit          string  `json:"unit"`
	LowStockAlert float64 `json:"low_stock_alert"`
}

type RestockProductDto struct {
	PackageCount   float64 `json:"package_count"`    // Cantidad de paquetes ingresados
	UnitPerPackage float64 `json:"unit_per_package"` // Cantidad por paquete
	Reason         string  `json:"reason"`           // Razón del movimiento (opcional)
	UnityPrice     float64 `json:"unity_price"`      // Precio unitario
}
