package repository

import "kiosco/models"

// ObtenerProductosActivos retorna todos los productos activos
func (r *Repositorio) ObtenerProductosActivos() ([]models.Producto, error) {
	query := `
		SELECT IdProducto, Nombre, PrecioUnitario, EstaActivo
		FROM Productos
		WHERE EstaActivo = true
		ORDER BY IdProducto
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo)
		if err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, nil
}

// ObtenerProductoPorId retorna un producto por su ID
func (r *Repositorio) ObtenerProductoPorId(idProducto int) (*models.Producto, error) {
	query := `
		SELECT IdProducto, Nombre, PrecioUnitario, EstaActivo
		FROM Productos
		WHERE IdProducto = $1
	`
	var p models.Producto
	err := r.db.QueryRow(query, idProducto).Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
