package repositories

import "kiosco/internal/models"

// ObtenerUsuarioPorNombre busca un usuario por su nombre de usuario
func (r *Repositorio) ObtenerUsuarioPorNombre(usuario string) (models.Usuario, error) {
	var u models.Usuario
	err := r.db.QueryRow(`
		SELECT id_usuario, usuario, contrasenha, puede_editar
		FROM usuarios WHERE usuario = ?
	`, usuario).Scan(&u.IdUsuario, &u.Usuario, &u.Contrasenha, &u.PuedeEditar)
	return u, err
}
