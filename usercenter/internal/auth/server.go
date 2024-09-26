package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"github.com/pengcainiao2/zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

type Server struct {
	//UnimplementedTokenServer

	conf *Config

	accessTokenProvider  TokenProvider
	refreshTokenProvider TokenProvider
}

var GlobalServer *Server

func New(conf *Config) *Server {
	//listen: :8080
	//logLevel: debug
	//
	//issuer: https://ttchat-api.zz7990.com
	//
	//accessToken:
	//keysFile: /etc/secrets/oauth/keys.json
	//keysReloadInterval: 1m
	//timeToLive: 6h
	//
	//refreshToken:
	//redis:
	//host: 10.0.1.30
	//port: 6379
	//timeToLive: 1440h
	//
	//jwkPublishAddr: :8081
	//jwkPublishPath: /.well-known/jwks.json
	conf.JwkPublishPath = "/.well-known/jwks.json"
	conf.JwkPublishAddr = ":8081"
	conf.Issuer = "http://penglonghui.cn"
	conf.AccessToken = &AccessTokenConfig{
		KeysFile:           "../keys.json",
		KeysReloadInterval: 1 * time.Minute,
		TimeToLive:         6 * time.Hour,
	}
	conf.RefreshToken = &RefreshTokenConfig{
		Redis: RedisConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
		TimeToLive: 1440 * time.Hour,
	}

	return &Server{conf: conf}
}

func (s *Server) Config() interface{} {
	return &s.conf
}

func (s *Server) OnConfigChange() {
	log.Infof("config updated: %v", s.conf)
}

func (s *Server) Init() error {
	logx.NewTraceLogger(context.Background()).Info().Msg(fmt.Sprintf("init begin: %v", s.conf))

	var err error

	s.refreshTokenProvider, err = NewRefreshTokenProvider(s.conf.RefreshToken)
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Msg(fmt.Sprintf("initNewRefreshTokenProvider fail, %v", err))

		return err
	}

	s.accessTokenProvider, err = NewAccessTokenProvider(s.conf.Issuer, s.conf.AccessToken)
	if err != nil {
		logx.NewTraceLogger(context.Background()).Err(err).Msg(fmt.Sprintf("NewAccessTokenProvider fail, %v", err))
		return err
	}

	go func() {
		err := serveJWKPublisher(s.conf.JwkPublishAddr, s.conf.JwkPublishPath, s)
		if err != nil {
			logx.NewTraceLogger(context.Background()).Err(err).Msg(fmt.Sprintf("serveJWKPublisher fail, %v", err))
		}
		log.Infof("serveJWKPublisher: %v", err)
	}()
	GlobalServer = s
	return nil
}

type GrantTokensRequest struct {
	TokenInfo *TokenInfo `protobuf:"bytes,1,opt,name=token_info,json=tokenInfo,proto3" json:"token_info,omitempty"`
}

type GrantTokensResponse struct {
	AccessToken  *GrantedToken `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken *GrantedToken `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

type GrantedToken struct {
	Type      TokenType `protobuf:"varint,1,opt,name=type,proto3,enum=azeroth.northrend.TokenType" json:"type,omitempty"`
	Token     string    `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresIn int64     `protobuf:"varint,3,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
	TokenId   string    `protobuf:"bytes,4,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
}

type TokenType int32

const (
	TokenType_UNDEFINED     TokenType = 0
	TokenType_ACCESS_TOKEN  TokenType = 1
	TokenType_REFRESH_TOKEN TokenType = 2
)

type TokenInfo struct {
	Uid         uint64   `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	DeviceId    string   `protobuf:"bytes,2,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	Scopes      []string `protobuf:"bytes,3,rep,name=scopes,proto3" json:"scopes,omitempty"`
	CountryCode string   `protobuf:"bytes,4,opt,name=country_code,json=countryCode,proto3" json:"country_code,omitempty"`
	IssueAt     uint64   `protobuf:"varint,5,opt,name=issue_at,json=issueAt,proto3" json:"issue_at,omitempty"`
}

func (s *Server) GrantTokens(ctx context.Context, in *GrantTokensRequest) (out *GrantTokensResponse, err error) {
	logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("GrantTokens, token info: (%v)", in.TokenInfo))

	out = &GrantTokensResponse{}
	if s.accessTokenProvider == nil {
		logx.NewTraceLogger(ctx).Info().Msg(fmt.Sprintf("555"))

	}
	out.AccessToken, err = s.accessTokenProvider.GrantToken(ctx, in.TokenInfo, s.conf.AccessToken.TimeToLive)
	if err != nil {
		return nil, err
	}
	out.RefreshToken, err = s.refreshTokenProvider.GrantToken(ctx, in.TokenInfo, s.conf.RefreshToken.TimeToLive)
	if err != nil {
		return nil, err
	}

	return
}

type ValidateTokenRequest struct {
	Type  TokenType `protobuf:"varint,1,opt,name=type,proto3,enum=azeroth.northrend.oauth.TokenType" json:"type,omitempty"`
	Token string    `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

type ValidateTokenResponse struct {
	TokenInfo *TokenInfo `protobuf:"bytes,1,opt,name=token_info,json=tokenInfo,proto3" json:"token_info,omitempty"` // indicates the request token is validated
}

func (s *Server) ValidateToken(ctx context.Context, in *ValidateTokenRequest) (out *ValidateTokenResponse,
	err error) {
	log.Debugf("ValidateToken, type: (%v), token: (%v)", in.Type, in.Token)

	out = &ValidateTokenResponse{}

	var p TokenProvider
	switch in.Type {
	case TokenType_ACCESS_TOKEN:
		p = s.accessTokenProvider
	case TokenType_REFRESH_TOKEN:
		p = s.refreshTokenProvider
	default:
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid token type: %s", in.Type)
	}

	out.TokenInfo, err = p.ValidateToken(ctx, in.Token)
	if err != nil {
		log.Errorf("ValidateToken %s %s %v", in.Type, in.Token, err)
	}
	return
}

type RevokeAllTokensRequest struct {
	Uid  uint64         `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Opts *RevokeOptions `protobuf:"bytes,2,opt,name=opts,proto3" json:"opts,omitempty"`
}

type RevokeOptions struct {
	ExcludeAccessTokens  []string `protobuf:"bytes,1,rep,name=exclude_access_tokens,json=excludeAccessTokens,proto3" json:"exclude_access_tokens,omitempty"`
	ExcludeRefreshTokens []string `protobuf:"bytes,3,rep,name=exclude_refresh_tokens,json=excludeRefreshTokens,proto3" json:"exclude_refresh_tokens,omitempty"`
	DeviceId             string   `protobuf:"bytes,4,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
}

type RevokeAllTokensResponse struct {
}

func (s *Server) RevokeAllTokens(ctx context.Context, in *RevokeAllTokensRequest) (*RevokeAllTokensResponse, error) {
	revoked, err := s.refreshTokenProvider.RevokeAllTokens(ctx, in.Uid, in.Opts.DeviceId, in.Opts.ExcludeRefreshTokens...)
	if err != nil {
		return nil, err
	}
	log.Infof("Revoke %d %v %v", in.Uid, in.Opts, revoked)
	return &RevokeAllTokensResponse{}, nil
}

func (s *Server) GetAccessTokenPublicKeys(ctx context.Context, in *GetAccessTokenPublicKeysRequest) (*GetAccessTokenPublicKeysResponse, error) {
	jwkSet := s.accessTokenProvider.GetPublicJWKSet(ctx, in.KeyId)
	for idx, k := range jwkSet.Keys {
		if !k.IsPublic() {
			jwkSet.Keys = append(jwkSet.Keys[0:idx], jwkSet.Keys[idx+1:]...)
		}
	}

	b, err := json.Marshal(jwkSet)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &GetAccessTokenPublicKeysResponse{JwkFormatKeys: string(b)}, nil
}
