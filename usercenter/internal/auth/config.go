package auth

type Config struct {
	Issuer       string
	AccessToken  *AccessTokenConfig
	RefreshToken *RefreshTokenConfig

	JwkPublishAddr string
	JwkPublishPath string
}
