package credential

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/util"
)

const (
	// accessTokenURL 获取 access_token 的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// stableAccessTokenURL 获取稳定版 access_token 的接口
	stableAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/stable_token"
	// workAccessTokenURL 企业微信获取 access_token 的接口
	workAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// CacheKeyOfficialAccountPrefix 微信公众号 cache key 前缀
	CacheKeyOfficialAccountPrefix = "gowechat_officialaccount_"
	// CacheKeyMiniProgramPrefix 小程序 cache key 前缀
	CacheKeyMiniProgramPrefix = "gowechat_miniprogram_"
	// CacheKeyWorkPrefix 企业微信 cache key 前缀
	CacheKeyWorkPrefix = "gowechat_work_"
)

// DefaultAccessToken 默认 AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewDefaultAccessToken new DefaultAccessToken
func NewDefaultAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache is ineed")
	}
	return &DefaultAccessToken{
		appID:           appID,
		appSecret:       appSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetAccessToken 获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *DefaultAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从 cache 中取
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.appID)

	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		var ok bool
		if accessToken, ok = val.(string); ok && accessToken != "" {
			return
		}
	}

	// 加上 lock，是为了防止在并发获取 token 时，cache 刚好失效，导致从微信服务器上获取到不同 token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()

	// 双检，防止重复从微信服务器获取
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		var ok bool
		if accessToken, ok = val.(string); ok && accessToken != "" {
			return
		}
	}

	// cache 失效，从微信服务器获取
	var resAccessToken ResAccessToken
	if resAccessToken, err = GetTokenFromServerContext(ctx, fmt.Sprintf(accessTokenURL, ak.appID, ak.appSecret)); err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 1500
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

	accessToken = resAccessToken.AccessToken
	return
}

// StableAccessToken 获取稳定版接口调用凭据 (与 getAccessToken 获取的调用凭证完全隔离，互不影响)
// 不强制更新 access_token，可用于不同环境不同服务而不需要分布式锁以及公用缓存，避免 access_token 争抢
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-access-token/getStableAccessToken.html
type StableAccessToken struct {
	appID          string
	appSecret      string
	cacheKeyPrefix string
	cache          cache.Cache
}

// NewStableAccessToken new StableAccessToken
func NewStableAccessToken(appID, appSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache is need")
	}
	return &StableAccessToken{
		appID:          appID,
		appSecret:      appSecret,
		cache:          cache,
		cacheKeyPrefix: cacheKeyPrefix,
	}
}

// GetAccessToken 获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *StableAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *StableAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 先从 cache 中取
	accessTokenCacheKey := fmt.Sprintf("%s_stable_access_token_%s", ak.cacheKeyPrefix, ak.appID)
	if val := ak.cache.Get(accessTokenCacheKey); val != nil {
		return val.(string), nil
	}

	// cache 失效，从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ak.GetAccessTokenDirectly(ctx, false)
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 300
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

	accessToken = resAccessToken.AccessToken
	return
}

// GetAccessTokenDirectly 从微信获取 access_token
func (ak *StableAccessToken) GetAccessTokenDirectly(ctx context.Context, forceRefresh bool) (resAccessToken ResAccessToken, err error) {
	b, err := util.PostJSONContext(ctx, stableAccessTokenURL, map[string]interface{}{
		"grant_type":    "client_credential",
		"appid":         ak.appID,
		"secret":        ak.appSecret,
		"force_refresh": forceRefresh,
	})
	if err != nil {
		return
	}

	if err = json.Unmarshal(b, &resAccessToken); err != nil {
		return
	}

	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get stable access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}
	return
}

// WorkAccessToken 企业微信 AccessToken 获取
type WorkAccessToken struct {
	CorpID          string
	CorpSecret      string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewWorkAccessToken new WorkAccessToken
func NewWorkAccessToken(corpID, corpSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenContextHandle {
	if cache == nil {
		panic("cache the not exist")
	}
	return &WorkAccessToken{
		CorpID:          corpID,
		CorpSecret:      corpSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// GetAccessToken 企业微信获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *WorkAccessToken) GetAccessToken() (accessToken string, err error) {
	return ak.GetAccessTokenContext(context.Background())
}

// GetAccessTokenContext 企业微信获取 access_token，先从 cache 中获取，没有则从服务端获取
func (ak *WorkAccessToken) GetAccessTokenContext(ctx context.Context) (accessToken string, err error) {
	// 加上 lock，是为了防止在并发获取 token 时，cache 刚好失效，导致从微信服务器上获取到不同 token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.CorpID)
	val := ak.cache.Get(accessTokenCacheKey)
	if val != nil {
		var ok bool
		if accessToken, ok = val.(string); !ok {
			accessToken = ""
		}
		return
	}

	// cache 失效，从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = GetTokenFromServerContext(ctx, fmt.Sprintf(workAccessTokenURL, ak.CorpID, ak.CorpSecret))
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 1500
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)

	accessToken = resAccessToken.AccessToken
	return
}

// GetTokenFromServer 强制从微信服务器获取 token
func GetTokenFromServer(url string) (resAccessToken ResAccessToken, err error) {
	return GetTokenFromServerContext(context.Background(), url)
}

// GetTokenFromServerContext 强制从微信服务器获取 token
func GetTokenFromServerContext(ctx context.Context, url string) (resAccessToken ResAccessToken, err error) {
	var body []byte
	body, err = util.HTTPGetContext(ctx, url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}
	return
}
