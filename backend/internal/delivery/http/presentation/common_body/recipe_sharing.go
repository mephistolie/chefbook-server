package common_body

type RecipeUserPublicKey struct {
	PublicKey string `json:"encrypted_public_key" binding:"required,min=10"`
}

type RecipeOwnerPrivateKey struct {
	PrivateKey string `json:"encrypted_private_key" binding:"required,min=10"`
}
