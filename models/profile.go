package models

// Profile kullanıcı profil modeli
type Profile struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	// Diğer alanlar
}

// ChangePasswordRequest şifre değiştirme isteği modeli
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
