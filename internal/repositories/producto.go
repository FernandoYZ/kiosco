package repositories

import "kiosco/internal/models"

// ObtenerTodosProductos retorna todos los productos (activos e inactivos)
func (r *Repositorio) ObtenerTodosProductos() ([]models.Producto, error) {
	rows, err := r.db.Query(`
		SELECT id_producto, nombre, precio_unitario, esta_activo
		FROM productos
		ORDER BY esta_activo DESC, id_producto
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		if err := rows.Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, rows.Err()
}

// InsertarProducto agrega un nuevo producto activo
func (r *Repositorio) InsertarProducto(nombre string, precio float64) (models.Producto, error) {
	result, err := r.db.Exec(`
		INSERT INTO productos (nombre, precio_unitario, esta_activo)
		VALUES (?, ?, 1)
	`, nombre, precio)
	if err != nil {
		return models.Producto{}, err
	}
	id, _ := result.LastInsertId()
	return models.Producto{
		IdProducto:     int(id),
		Nombre:         nombre,
		PrecioUnitario: precio,
		EstaActivo:     true,
	}, nil
}

// ActualizarProducto modifica nombre y precio de un producto
func (r *Repositorio) ActualizarProducto(id int, nombre string, precio float64) error {
	_, err := r.db.Exec(`
		UPDATE productos SET nombre = ?, precio_unitario = ? WHERE id_producto = ?
	`, nombre, precio, id)
	return err
}

// CambiarEstadoProducto habilita o deshabilita un producto
func (r *Repositorio) CambiarEstadoProducto(id int, activo bool) error {
	estado := 0
	if activo {
		estado = 1
	}
	_, err := r.db.Exec(`
		UPDATE productos SET esta_activo = ? WHERE id_producto = ?
	`, estado, id)
	return err
}

// ObtenerProductosActivos retorna todos los productos activos
func (r *Repositorio) ObtenerProductosActivos() ([]models.Producto, error) {
	rows, err := r.db.Query(`
		SELECT id_producto, nombre, precio_unitario, esta_activo
		FROM productos
		WHERE esta_activo = 1
		ORDER BY id_producto
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		if err := rows.Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}

	return productos, rows.Err()
}

// ObtenerProductoPorId retorna un producto por su ID
func (r *Repositorio) ObtenerProductoPorId(idProducto int) (*models.Producto, error) {
	var p models.Producto
	err := r.db.QueryRow(`
		SELECT id_producto, nombre, precio_unitario, esta_activo
		FROM productos WHERE id_producto = ?
	`, idProducto).Scan(&p.IdProducto, &p.Nombre, &p.PrecioUnitario, &p.EstaActivo)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
