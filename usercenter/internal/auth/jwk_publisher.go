package auth

import (
	"github.com/gin-gonic/gin"
)

type GetAccessTokenPublicKeysRequest struct {
	// OPTIONAL, if not specified, server returns all keys
	KeyId string `protobuf:"bytes,1,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
}

type GetAccessTokenPublicKeysResponse struct {
	JwkFormatKeys string `protobuf:"bytes,1,opt,name=jwk_format_keys,json=jwkFormatKeys,proto3" json:"jwk_format_keys,omitempty"`
}

func serveJWKPublisher(addr, path string, s *Server) error {
	engine := gin.New()

	engine.Use(gin.Logger(), gin.Recovery())
	engine.GET(path, func(c *gin.Context) {
		r, err := s.GetAccessTokenPublicKeys(c.Request.Context(), &GetAccessTokenPublicKeysRequest{})
		if err == nil {
			c.String(200, r.JwkFormatKeys)
		} else {
			c.JSON(500, nil)
		}
	})

	return engine.Run(addr)
}
