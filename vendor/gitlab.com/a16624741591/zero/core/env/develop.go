package env

import (
	"os"
	"strings"
)

const (
	aliResourceAccessKey        = "ALIRES_ACCESS_KEY"
	aliResourceAccessSecret     = "ALIRES_ACCESS_SECRET"
	ginMode                     = "GIN_MODE"
	ossUploadCallback           = "OSS_UPLOAD_CALLBACK"
	ossBucketHost               = "OSS_BUCKET_HOST"
	ossEndpoint                 = "OSS_ENDPOINT"
	ossBucketName               = "BUCKET_NAME"
	dbDSN                       = "DB_DSN"
	redisAddr                   = "REDIS_ADDR"
	redisPwd                    = "REDIS_PWD"
	redisMasterName             = "REDIS_MASTER_NAME"
	wechatAppID                 = "WECHAT_APP_ID"
	wechatAppSecret             = "WECHAT_APP_SECRET"
	wechatOfficialID            = "WECHAT_OFFICIAL_ID"
	wechatOfficialSecret        = "WECHAT_OFFICIAL_SECRET"
	releaseMode                 = "RELEASE_MODE"
	rateBurst                   = "RATE_BURST"
	secureRefer                 = "SECURE_REFER"
	sensorsAddr                 = "SA_URL"
	officialAccountName         = "OFFICIAL_ACCOUNT_NAME"
	wechatOpenAppID             = "WECHAT_OPEN_APP_ID"
	wechatOpenAppSecret         = "WECHAT_OPEN_APP_SECRET"
	appPackageName              = "APP_PACKAGE_NAME"
	appleP8CertificateKey       = "APPLE_P8_KEY"
	appleP8Certificate          = "APPLE_P8_FILE"
	appleTeamID                 = "APPLE_TEAM_ID"
	domainName                  = "DOMAIN_NAME"
	tpnsAndroidAccessID         = "TPNS_ANDROID_ACCESS_ID"
	tpnsAndroidSecretKey        = "TPNS_ANDROID_SECRET_KEY"
	tpnsIOSAccessID             = "TPNS_IOS_ACCESS_ID"
	tpnsIOSSecretKey            = "TPNS_IOS_SECRET_KEY"
	workWechatMerchantSecretKey = "WORK_WECHAT_MERCHANT_SECRET_KEY"
	workWechatMerchantID        = "WORK_WECHAT_MERCHANT_ID"
	newTaskInviteHref           = "NEW_TASK_INVITE_HREF"
	newProjectInviteHref        = "NEW_PROJECT_INVITE_HREF"
	newMeetingInviteHref        = "NEW_MEETING_INVITE_HREF"
	newWorkspaceInviteHref      = "NEW_WORKSPACE_INVITE_HREF"
	newBindingInviteHref        = "NEW_BINDING_INVITE_HREF"
	getRobotTokenUrl            = "GET_ROBOT_TOKEN_URL"
	getRobotTokenAuthorization  = "GET_ROBOT_TOKEN_AUTHORIZATION"
	sendMessageToWechatUrl      = "SEND_MESSAGE_TO_WECHAT_URL"
	downloadFileUrl             = "DOWNLOAD_FILE_URL"
	nsqAddr                     = "NSQ_ADDR"
	etcdAddr                    = "ETCD_ADDR"
	namespace                   = "NAMESPACE"
	jaegerCollectorAddr         = "JAEGER_COLLECTOR_ADDR"
	zincHost                    = "ZINC_HOST"
	zincUser                    = "ZINC_USER"
	zincPassword                = "ZINC_PASSWORD"
	zincAddress                 = "ZINC_ADDR"
	elasticAddr                 = "ELASTIC_ADDR"
	elasticUser                 = "ELASTIC_USER"
	elasticPassword             = "ELASTIC_PASSWORD"
	privateUploadPath           = "PRIVATE_UPLOAD_PATH"
	privateAuthCode             = "PRIVATE_AUTH_CODE"
	privateAuthMode             = "PRIVATE_AUTH_MODE"
	privateAuthSSHHOST          = "PRIVATE_AUTH_SSH_HOST"
	privateAuthSSHUser          = "PRIVATE_AUTH_SSH_USER"
	privateAuthSSHKeyFile       = "PRIVATE_AUTH_SSH_KEY_FILE"
	privateAuthHTTPHOST         = "PRIVATE_AUTH_HTTP_HOST"
	privateStaticHOST           = "PRIVATE_STATIC_HOST"
	privateWebserviceHOST       = "PRIVATE_WEBSERVICE_HOST"
	wechatMchID                 = "WECHAT_MCH_ID"
	wechatAPIv3Key              = "WECHAT_API_V3_KRY"
	wechatPaySerialNo           = "WECHAT_SERIAL_NO"
	wechatPayPrivateKeyContent  = "WECHAT_PAY_PRIVATE_KEY_CONTENT"
	yzjAppID                    = "YZJ_APP_ID"
	yzjAppSecret                = "YZJ_APP_SECRET"
	yzjSignKey                  = "YZJ_SIGN_KEY"
	yzjEid                      = "YZJ_EID"
	yzjPub                      = "YZJ_PUB"
	yzjPubKey                   = "YZJ_PUB_KEY"
)

var (
	configMap = envConfig{}
)
var (
	GinMode                = configMap.getValue(ginMode)
	Namespace              = configMap.getValue(namespace)
	NSQAddress             = configMap.getValue(nsqAddr)
	ETCDAddress            = configMap.getValue(etcdAddr)
	JaegerCollectorAddress = configMap.getValue(jaegerCollectorAddr)
	ZincAddress            = configMap.getValue(zincAddress)
	DomainName             = configMap.getValue(domainName)
	//AppPackageName 应用包名
	AppPackageName = configMap.getValue(appPackageName)
	//AppleP8CertificateKey 苹果APP在使用苹果账号授权时所需要的p8 证书key
	AppleP8CertificateKey = configMap.getValue(appleP8CertificateKey)
	AppleP8Certificate    = configMap.getValue(appleP8Certificate)
	AppleTeamID           = configMap.getValue(appleTeamID)
	// 用于访问阿里云资源ID
	AliResourcesAccessKey = configMap.getValue(aliResourceAccessKey)
	// 用于访问阿里云资源秘钥
	AliResourcesAccessSecret = configMap.getValue(aliResourceAccessSecret)
	// OSS文件上传成功后的回调地址
	OssUploadCallback = configMap.getValue(ossUploadCallback)
	// 所访问bucket的地址
	OssBucketHost = configMap.getValue(ossBucketHost)
	// OSS端点地址
	OssEndpoint = configMap.getValue(ossEndpoint)
	// bucket名称
	OssBucketName = configMap.getValue(ossBucketName)
	// 数据库地址
	DbDSN = configMap.getValue(dbDSN)
	// Redis地址
	RedisAddr = configMap.getValue(redisAddr)
	// Redis密码
	RedisPwd = configMap.getValue(redisPwd)
	// RedisMasterName 哨兵模式master_name
	RedisMasterName = configMap.getValue(redisMasterName)
	// 微信小程序ID
	WechatAppID = configMap.getValue(wechatAppID)
	// 微信小程序秘钥
	WechatAppSecret = configMap.getValue(wechatAppSecret)
	// 微信公众号ID
	WechatOfficialID = configMap.getValue(wechatOfficialID)
	// 微信公众号秘钥
	WechatOfficialSecret = configMap.getValue(wechatOfficialSecret)
	// 发布模式
	ReleaseMode = configMap.getValue(releaseMode)
	// 限流
	RateBurst = configMap.getValue(rateBurst)

	//Deprecated 仅usercenter使用，使用grpc接口代替
	SecureRefer = configMap.getValue(secureRefer)
	// 神策地址
	SensorsAddr = configMap.getValue(sensorsAddr)
	//Deprecated
	OfficialAccountName = configMap.getValue(officialAccountName)
	// 微信开发平台应用ID（移动端第三方授权登录）
	WechatOpenAppID = configMap.getValue(wechatOpenAppID)
	// 微信开发平台应用秘钥
	WechatOpenAppSecret = configMap.getValue(wechatOpenAppSecret)
	// 腾讯推送配置
	TpnsIOSAccessId      = configMap.getValue(tpnsIOSAccessID)
	TpnsIOSSecretKey     = configMap.getValue(tpnsIOSSecretKey)
	TpnsAndroidAccessId  = configMap.getValue(tpnsAndroidAccessID)
	TpnsAndroidSecretKey = configMap.getValue(tpnsAndroidSecretKey)
	// 企业微信机器人密钥
	WorkWechatMerchantSecretKey = configMap.getValue(workWechatMerchantSecretKey)
	// 企业微信机器人商家ID
	WorkWechatMerchantID = configMap.getValue(workWechatMerchantID)
	// 微信机器人事项邀请链接
	NewTaskInviteHref = configMap.getValue(newTaskInviteHref)
	// 微信机器人项目邀请链接
	NewProjectInviteHref = configMap.getValue(newProjectInviteHref)
	// 微信机器人会议邀请链接
	NewMeetingInviteHref = configMap.getValue(newMeetingInviteHref)
	// 微信机器人空间邀请链接
	NewWorkspaceInviteHref = configMap.getValue(newWorkspaceInviteHref)
	// 微信机器人绑定链接
	NewBindingInviteHref = configMap.getValue(newBindingInviteHref)
	// 微信机器人获取token链接
	GetRobotTokenUrl = configMap.getValue(getRobotTokenUrl)
	// 微信机器人获取token证书
	GetRobotTokenAuthorization = configMap.getValue(getRobotTokenAuthorization)
	// 微信机器人发送微信消息链接
	SendMessageToWechatUrl = configMap.getValue(sendMessageToWechatUrl)
	// 微信机器人下载文件链接
	DownloadFileUrl = configMap.getValue(downloadFileUrl)

	// zincsearch相关
	ZincHost     = configMap.getValue(zincHost)
	ZincUser     = configMap.getValue(zincUser)
	ZincPassword = configMap.getValue(zincPassword)

	// elastic相关
	ElasticAddr     = configMap.getValue(elasticAddr)
	ElasticUser     = configMap.getValue(elasticUser)
	ElasticPassword = configMap.getValue(elasticPassword)

	mpAppSecretMappingMap = map[string]string{
		"wxfca9e569e0d02bb2": "57014ca43fcfe830f6702fe7fb137268", //飞项小助手
		"wx47498aca992164df": "f292544ca09b5e462a84060387c5b880", //飞项|事项协作..
		"wxf4725ab51e893fad": "ddf5a3f4f8223d29642714b6d055c727", //任务派发
		"wxefb12780cb146b1e": "f43abcfdab92d3c60761fbc9adef0202", //清单待办
		"wx7a2848ce31003e27": "abfd1693671da04d063fcc94d97d0f7a", //事项协作
		"wx7773e2a7d0fb41c0": "0904da4e393e4ece022541cce548c630", //清单日程
		"wx5c12b325da25825f": "e1dc01b481d60beb2a473655a057e607", //待办协作
		"wx738e469c5ae41835": "6a8a005e6b1dc98e6402811fe582ce90", //群协作
		"wx3501d4809a0de8d7": "f7bdef13588215e9a32eac3c5fd8fa98", //群项目
		"wx2ed653ad5f070600": "415eda9c87dd731ca0fc1fb327c9571d", //群事项
		"wxfbf5f9ef7b3410f1": "a0ee73239284e27f53b1484a6d3f758f", //群发待办
	}
	// 私有化文件上传路径
	PrivateUploadPath = configMap.getValue(privateUploadPath)
	// PrivateAuthCode 授权码
	PrivateAuthCode = configMap.getValue(privateAuthCode)
	// PrivateAuthMode 私有化授权模式
	PrivateAuthMode = configMap.getValue(privateAuthMode)
	// PrivateAuthSSHHOST 私有化授权ssh地址
	PrivateAuthSSHHOST = configMap.getValue(privateAuthSSHHOST)
	// PrivateAuthSSHUser 私有化授权ssh用户
	PrivateAuthSSHUser = configMap.getValue(privateAuthSSHUser)
	// PrivateAuthSSHKeyFile 私有化授权ssh公钥地址
	PrivateAuthSSHKeyFile = configMap.getValue(privateAuthSSHKeyFile)
	// PrivateAuthHTTPHOST 私有化授权http地址
	PrivateAuthHTTPHOST = configMap.getValue(privateAuthHTTPHOST)
	// PrivateStaticHOST 私有化静态文件 host 地址
	PrivateStaticHOST = configMap.getValue(privateStaticHOST)
	// PrivateWebserviceHOST 私有化webservice地址
	PrivateWebserviceHOST = configMap.getValue(privateWebserviceHOST)

	// WechatMchID 微信支付商户ID
	WechatMchID                = configMap.getValue(wechatMchID)
	WechatAPIv3Key             = configMap.getValue(wechatAPIv3Key)
	WechatPaySerialNo          = configMap.getValue(wechatPaySerialNo)
	WechatPayPrivateKeyContent = configMap.getValue(wechatPayPrivateKeyContent)
	// YzjAppID 云之家appid
	YzjAppID = configMap.getValue(yzjAppID)
	// YzjAppSecret 云之家app_secret
	YzjAppSecret = configMap.getValue(yzjAppSecret)
	// YzjSignKey 云之家签名key
	YzjSignKey = configMap.getValue(yzjSignKey)
	// YzjEid 云之家团队id
	YzjEid = configMap.getValue(yzjEid)
	// YzjPub 云之家公共号账号
	YzjPub = configMap.getValue(yzjPub)
	// YzjPubKey 云之家公共号秘钥
	YzjPubKey = configMap.getValue(yzjPubKey)
)

type envConfig map[string]string

func (env envConfig) getValue(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if v, ok := env[key]; ok {
		return v
	}
	return ""
}

//IsProduction 是否为生产环境
func IsProduction() bool {
	releaseMode := configMap.getValue(releaseMode)
	return configMap != nil && (releaseMode == "production" || releaseMode == "pre-production") && IsRunningInK8s()
}

//IsRunningInK8s 服务是否运行在k8s环境中
func IsRunningInK8s() bool {
	var v = os.Getenv("KUBERNETES_PORT")
	return v != ""
}

//IsRunningInAliyun 服务是否运行在aliyun
func IsRunningInAliyun() bool {
	return strings.HasSuffix(configMap.getValue(domainName), "flyele.net")
}

func IsDevMode() bool {
	return !IsRunningInK8s() && (ReleaseMode == "develop" || ReleaseMode == "")
}

// IsPrivate 是否是私有化环境
func IsPrivate() bool {
	return configMap.getValue(releaseMode) == "private"
}

func GetAppMapping() map[string]string {
	return mpAppSecretMappingMap
}

//EnableTracing 是否启用opentracing
func EnableTracing() bool {
	//return false
	return IsRunningInK8s() //&& ReleaseMode == "release"
}
