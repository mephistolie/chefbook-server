package request_body

type PasswordChanging struct {
	OldPassword string `json:"old_password" binding:"max=64"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=64"`
}

func (p *PasswordChanging) Validate() error {
	return validatePassword(p.NewPassword)
}

type Username struct {
	Username *string `json:"username" binding:"max=40"`
}
