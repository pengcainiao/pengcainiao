package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/martian/log"
	"gopkg.in/square/go-jose.v2"
	"os"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	minimumKeysReloadInterval = time.Second * 10
)

type AccessTokenConfig struct {
	KeysFile           string
	KeysReloadInterval time.Duration
	TimeToLive         time.Duration
}

type accessTokenClaims struct {
	jwt.Claims
	DeviceID    string   `json:"dev"`
	CountryCode string   `json:"ctc"`
	Scopes      []string `json:"scopes"`
}

type accessTokenProvider struct {
	conf *AccessTokenConfig

	iss          string
	signKeys     jose.JSONWebKeySet
	verifyKeySet jose.JSONWebKeySet

	idgen *sonyflake.Sonyflake
}

func NewAccessTokenProvider(iss string, conf *AccessTokenConfig) (TokenProvider, error) {
	p := &accessTokenProvider{
		conf:  conf,
		iss:   iss,
		idgen: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}

	err := p.reloadKeys()
	if err != nil {
		log.Errorf("reload keys: %v", err)
		return nil, err
	}

	go func() {
		interval := conf.KeysReloadInterval
		if interval < minimumKeysReloadInterval {
			interval = minimumKeysReloadInterval
		}
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				err := p.reloadKeys()
				if err != nil {
					log.Errorf("reload keys: %v", err)
				}
			}
		}
	}()

	return p, nil
}

func (p *accessTokenProvider) reloadKeys() error {
	jwkSet, err := loadKeys(p.conf.KeysFile)
	if err != nil {
		return err
	}
	if len(jwkSet.Keys) == 0 {
		return errors.New("empty jwk set")
	}

	publicKeys := make([]jose.JSONWebKey, 0, len(jwkSet.Keys))
	for _, k := range jwkSet.Keys {
		pk := k.Public()
		publicKeys = append(publicKeys, pk)
	}

	p.signKeys, p.verifyKeySet = *jwkSet, jose.JSONWebKeySet{Keys: publicKeys}

	return nil
}

func (p *accessTokenProvider) GrantToken(ctx context.Context, tokenInfo *TokenInfo, ttl time.Duration) (*GrantedToken, error) {
	log.Debugf("grant token: %d %s %s", tokenInfo.Uid, tokenInfo.DeviceId, ttl)
	now := time.Now()

	id, err := p.idgen.NextID()
	if err != nil {
		return nil, err
	}

	tokenId := strconv.FormatUint(id, 36)
	claims := &accessTokenClaims{
		Claims: jwt.Claims{
			Issuer:   p.iss,
			ID:       tokenId,
			Expiry:   jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt: jwt.NewNumericDate(now),
			Subject:  strconv.FormatUint(tokenInfo.Uid, 10),
		},
		DeviceID:    tokenInfo.DeviceId,
		Scopes:      tokenInfo.Scopes,
		CountryCode: tokenInfo.CountryCode,
	}

	signKey := p.signKeys.Keys[id%uint64(len(p.signKeys.Keys))]

	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(signKey.Algorithm),
		Key:       signKey,
	}, nil)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		return nil, err
	}

	return &GrantedToken{
		Type:      TokenType_ACCESS_TOKEN,
		Token:     token,
		TokenId:   tokenId,
		ExpiresIn: int64(ttl.Seconds()),
	}, nil
}

func (p *accessTokenProvider) ValidateToken(ctx context.Context, encodedToken string) (*TokenInfo, error) {
	var claims accessTokenClaims

	token, err := jwt.ParseSigned(encodedToken)
	if err != nil {
		return nil, invalidTokenError(err)
	}

	err = token.Claims(&p.verifyKeySet, &claims)
	if err != nil {
		log.Errorf("parse token %s: %v", encodedToken, err)
		return nil, invalidTokenError(err)
	}

	err = claims.Validate(jwt.Expected{Time: time.Now()})
	if err != nil {
		log.Errorf("parse token %s: %v", encodedToken, err)
		return nil, invalidTokenError(err)
	}

	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, invalidTokenError(errInvalidSubject)
	}

	return &TokenInfo{Uid: uid, DeviceId: claims.DeviceID}, nil
}

func (p *accessTokenProvider) RevokeToken(ctx context.Context, token string, def time.Duration) error {
	return errAccessTokenNotRevocable
}

func (p *accessTokenProvider) RevokeAllTokens(ctx context.Context, uid uint64, deviceId string, excludes ...string) ([]string, error) {
	return nil, errAccessTokenNotRevocable
}

func (p *accessTokenProvider) GetPublicJWKSet(ctx context.Context, keyID string) jose.JSONWebKeySet {
	if keyID == "" {
		return p.verifyKeySet
	}
	return jose.JSONWebKeySet{Keys: p.verifyKeySet.Key(keyID)}
}

func loadKeys(path string) (*jose.JSONWebKeySet, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var jwkSet jose.JSONWebKeySet
	err = json.NewDecoder(f).Decode(&jwkSet)
	if err != nil {
		return nil, err
	}
	return &jwkSet, nil
}
