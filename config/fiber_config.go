package config

import (
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis"
	"github.com/spf13/viper"
)

const (
	KByte int = 1024
	MByte int = KByte * KByte
)

func FiberConfigSetDefault() {
	_ = viper.BindEnv("fiber.prefork", "FIBER_PREFORK")
	viper.SetDefault("fiber.prefork", false)
	_ = viper.BindEnv("fiber.serverHeader", "FIBER_SERVER_HEADER")
	viper.SetDefault("fiber.serverHeader", "")
	_ = viper.BindEnv("fiber.strictRouting", "FIBER_STRICT_ROUTING")
	viper.SetDefault("fiber.strictRouting", true)
	_ = viper.BindEnv("fiber.caseSensitive", "FIBER_CASE_SENSITIVE")
	viper.SetDefault("fiber.caseSensitive", true)
	_ = viper.BindEnv("fiber.immutable", "FIBER_IMMUTABLE")
	viper.SetDefault("fiber.immutable", false)
	_ = viper.BindEnv("fiber.unescapePath", "FIBER_UNESCAPE_PATH")
	viper.SetDefault("fiber.unescapePath", false)
	_ = viper.BindEnv("fiber.etag", "FIBER_ETAG")
	viper.SetDefault("fiber.etag", true)
	_ = viper.BindEnv("fiber.bodyLimit", "FIBER_BODY_LIMIT")
	viper.SetDefault("fiber.bodyLimit", 4*MByte)
	_ = viper.BindEnv("fiber.concurrency", "FIBER_CONCURRENCY")
	viper.SetDefault("fiber.concurrency", 256*1024)
	_ = viper.BindEnv("fiber.readTimeout", "FIBER_READ_TIMEOUT")
	viper.SetDefault("fiber.readTimeout", 4*time.Second)
	_ = viper.BindEnv("fiber.writeTimeout", "FIBER_WRITE_TIMEOUT")
	viper.SetDefault("fiber.writeTimeout", 4*time.Second)
	_ = viper.BindEnv("fiber.idleTimeout", "FIBER_IDLE_TIMEOUT")
	viper.SetDefault("fiber.idleTimeout", 4*time.Second)
	_ = viper.BindEnv("fiber.readBufferSize", "FIBER_READ_BUFFER_SIZE")
	viper.SetDefault("fiber.readBufferSize", 4*KByte)
	_ = viper.BindEnv("fiber.writeBufferSize", "FIBER_WRITE_BUFFER_SIZE")
	viper.SetDefault("fiber.writeBufferSize", 4*KByte)
	_ = viper.BindEnv("fiber.compressedFileSuffix", "FIBER_COMPRESSED_FILE_SUFFIX")
	viper.SetDefault("fiber.compressedFileSuffix", ".fiber.gz")
	_ = viper.BindEnv("fiber.proxyHeader", "FIBER_PROXY_HEADER")
	viper.SetDefault("fiber.proxyHeader", "")
	_ = viper.BindEnv("fiber.getOnly", "FIBER_GET_ONLY")
	viper.SetDefault("fiber.getOnly", false)
	_ = viper.BindEnv("fiber.disableKeepalive", "FIBER_DISABLE_KEEPALIVE")
	viper.SetDefault("fiber.disableKeepalive", false)
	_ = viper.BindEnv("fiber.disableDefaultDate", "FIBER_DISABLE_DEFAULT_DATE")
	viper.SetDefault("fiber.disableDefaultDate", false)
	_ = viper.BindEnv("fiber.disableDefaultContentType", "FIBER_DISABLE_DEFAULT_CONTENT_TYPE")
	viper.SetDefault("fiber.disableDefaultContentType", false)
	_ = viper.BindEnv("fiber.disableHeaderNormalizing", "FIBER_DISABLE_HEADER_NORMALIZING")
	viper.SetDefault("fiber.disableHeaderNormalizing", false)
	_ = viper.BindEnv("fiber.disableStartupMessage", "FIBER_DISABLE_STARTUP_MESSAGE")
	viper.SetDefault("fiber.disableStartupMessage", false)
	_ = viper.BindEnv("fiber.appName", "FIBER_APP_NAME")
	viper.SetDefault("fiber.appName", "")
	_ = viper.BindEnv("fiber.streamRequestBody", "FIBER_STREAM_REQUEST_BODY")
	viper.SetDefault("fiber.streamRequestBody", true)
	_ = viper.BindEnv("fiber.disablePreParseMultipartForm", "FIBER_DISABLE_PRE_PARSE_MULTIPART_FORM")
	viper.SetDefault("fiber.disablePreParseMultipartForm", false)
	_ = viper.BindEnv("fiber.reduceMemoryUsage", "FIBER_REDUCE_MEMORY_USAGE")
	viper.SetDefault("fiber.reduceMemoryUsage", false)
	_ = viper.BindEnv("fiber.network", "FIBER_NETWORK")
	viper.SetDefault("fiber.network", fiber.NetworkTCP4)
	_ = viper.BindEnv("fiber.enableTrustedProxyCheck", "FIBER_ENABLE_TRUSTED_PROXY_CHECK")
	viper.SetDefault("fiber.enableTrustedProxyCheck", false)
	_ = viper.BindEnv("fiber.trustedProxies", "FIBER_TRUSTED_PROXIES")
	viper.SetDefault("fiber.trustedProxies", []string{})
	_ = viper.BindEnv("fiber.enableIpValidation", "FIBER_ENABLE_IP_VALIDATION")
	viper.SetDefault("fiber.enableIpValidation", true)
	_ = viper.BindEnv("fiber.enablePrintRoutes", "FIBER_ENABLE_PRINT_ROUTES")
	viper.SetDefault("fiber.enablePrintRoutes", false)

	_ = viper.BindEnv("fiber.encryptCookie.key", "FIBER_ENCRYPT_COOKIE_KEY")
	viper.SetDefault("fiber.encryptCookie.key", "secret-thirty-2-character-string")

	_ = viper.BindEnv("fiber.csrf.enable", "FIBER_CSRF_ENABLE")
	viper.SetDefault("fiber.csrf.enable", false)
	_ = viper.BindEnv("fiber.csrf.cookieSecure", "FIBER_CSRF_COOKIE_SECURE")
	viper.SetDefault("fiber.csrf.cookieSecure", true)
	_ = viper.BindEnv("fiber.csrf.cookieHttpOnly", "FIBER_CSRF_COOKIE_HTTP_ONLY")
	viper.SetDefault("fiber.csrf.cookieHttpOnly", true)

	_ = viper.BindEnv("fiber.cache.expiration", "FIBER_CACHE_EXPIRATION")
	viper.SetDefault("fiber.cache.expiration", 1*time.Minute)
	_ = viper.BindEnv("fiber.cache.control", "FIBER_CACHE_CONTROL")
	viper.SetDefault("fiber.cache.control", true)
	_ = viper.BindEnv("fiber.cache.header", "FIBER_CACHE_HEADER")
	viper.SetDefault("fiber.cache.header", "X-Cache")

	_ = viper.BindEnv("fiber.limiter.max", "FIBER_LIMITER_MAX")
	viper.SetDefault("fiber.limiter.max", 5)
	_ = viper.BindEnv("fiber.limiter.expiration", "FIBER_LIMITER_EXPIRATION")
	viper.SetDefault("fiber.limiter.expiration", 1*time.Minute)
	_ = viper.BindEnv("fiber.limiter.skipFailedRequests", "FIBER_LIMITER_SKIP_FAILED_REQUESTS")
	viper.SetDefault("fiber.limiter.skipFailedRequests", false)
	_ = viper.BindEnv("fiber.limiter.skipSuccessfulRequests", "FIBER_LIMITER_SKIP_SUCCESSFUL_REQUESTS")
	viper.SetDefault("fiber.limiter.skipSuccessfulRequests", false)

	_ = viper.BindEnv("fiber.cors.allowOrigins", "FIBER_CORS_ALLOW_ORIGINS")
	viper.SetDefault("fiber.cors.allowOrigins", "*")
	_ = viper.BindEnv("fiber.cors.allowMethods", "FIBER_CORS_ALLOW_METHODS")
	viper.SetDefault("fiber.cors.allowMethods", "GET POST HEAD PUT DELETE PATCH")
	_ = viper.BindEnv("fiber.cors.allowHeaders", "FIBER_CORS_ALLOW_HEADERS")
	viper.SetDefault("fiber.cors.allowHeaders", "")
	_ = viper.BindEnv("fiber.cors.exposeHeaders", "FIBER_CORS_EXPOSE_HEADERS")
	viper.SetDefault("fiber.cors.exposeHeaders", "")
	_ = viper.BindEnv("fiber.cors.allowCredentials", "FIBER_CORS_ALLOW_CREDENTIALS")
	viper.SetDefault("fiber.cors.allowCredentials", true)
	_ = viper.BindEnv("fiber.cors.maxAge", "FIBER_CORS_MAX_AGE")
	viper.SetDefault("fiber.cors.maxAge", 0)
}

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:       viper.GetBool("fiber.prefork"),
		ServerHeader:  viper.GetString("fiber.serverHeader"),
		StrictRouting: viper.GetBool("fiber.strictRouting"),
		CaseSensitive: viper.GetBool("fiber.caseSensitive"),
		Immutable:     viper.GetBool("fiber.immutable"),
		UnescapePath:  viper.GetBool("fiber.unescapePath"),
		// ETag:                 false,
		BodyLimit:   viper.GetInt("fiber.bodyLimit"),
		Concurrency: viper.GetInt("fiber.concurrency"),
		// Views:                nil,
		// ViewsLayout:          "",
		// PassLocalsToViews:    false,
		ReadTimeout:          viper.GetDuration("fiber.readTimeout"),
		WriteTimeout:         viper.GetDuration("fiber.writeTimeout"),
		IdleTimeout:          viper.GetDuration("fiber.idleTimeout"),
		ReadBufferSize:       viper.GetInt("fiber.readBufferSize"),
		WriteBufferSize:      viper.GetInt("fiber.writeBufferSize"),
		CompressedFileSuffix: viper.GetString("fiber.compressedFileSuffix"),
		ProxyHeader:          viper.GetString("fiber.proxyHeader"),
		GETOnly:              viper.GetBool("fiber.getOnly"),
		// ErrorHandler: func(*fiber.Ctx, error) error {
		// },
		DisableKeepalive:             viper.GetBool("fiber.disableKeepalive"),
		DisableDefaultDate:           viper.GetBool("fiber.disableDefaultDate"),
		DisableDefaultContentType:    viper.GetBool("fiber.disableDefaultContentType"),
		DisableHeaderNormalizing:     viper.GetBool("fiber.disableHeaderNormalizing"),
		DisableStartupMessage:        viper.GetBool("fiber.disableStartupMessage"),
		AppName:                      viper.GetString("fiber.appName"),
		StreamRequestBody:            viper.GetBool("fiber.streamRequestBody"),
		DisablePreParseMultipartForm: viper.GetBool("fiber.disablePreParseMultipartForm"),
		ReduceMemoryUsage:            viper.GetBool("fiber.reduceMemoryUsage"),
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		// XMLEncoder: func(v interface{}) ([]byte, error) {
		// },
		Network:                 viper.GetString("fiber.network"),
		EnableTrustedProxyCheck: viper.GetBool("fiber.enableTrustedProxyCheck"),
		TrustedProxies:          viper.GetStringSlice("fiber.trustedProxies"),
		EnableIPValidation:      viper.GetBool("fiber.enableIpValidation"),
		EnablePrintRoutes:       viper.GetBool("fiber.enablePrintRoutes"),
		// ColorScheme:             fiber.Colors{},
		// RequestMethods:          []string{},
	}
}

func NewFiberRedisStorageConfig() redis.Config {
	c := redis.Config{}
	if viper.GetBool("redis.auth") {
		c.Username = viper.GetString("redis.username")
		c.Password = viper.GetString("redis.password")
	}
	c.Host = viper.GetString("redis.host")
	c.Port = viper.GetInt("redis.port")
	c.Database = viper.GetInt("redis.database")
	c.PoolSize = viper.GetInt("redis.poolSize")
	return c
}

func NewFiberCacheConfig(s fiber.Storage) cache.Config {
	return cache.Config{
		Expiration:   viper.GetDuration("fiber.cache.expiration"),
		CacheHeader:  viper.GetString("fiber.cache.header"),
		CacheControl: viper.GetBool("fiber.cache.control"),
		Storage:      s,
	}
}

func NewFiberLimiterConfig(s fiber.Storage) limiter.Config {
	return limiter.Config{
		Max:                    viper.GetInt("fiber.limiter.max"),
		Expiration:             viper.GetDuration("fiber.limiter.expiration"),
		SkipFailedRequests:     viper.GetBool("fiber.limiter.skipFailedRequests"),
		SkipSuccessfulRequests: viper.GetBool("fiber.limiter.skipSuccessfulRequests"),
		Storage:                s,
		LimiterMiddleware:      limiter.SlidingWindow{},
		Next: func(c *fiber.Ctx) bool {
			switch c.IP() {
			case "127.0.0.1":
				return true
			default:
				return false
			}
		},
	}
}

func NewFiberCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     strings.Join(viper.GetStringSlice("fiber.cors.allowOrigins"), ","),
		AllowMethods:     strings.Join(viper.GetStringSlice("fiber.cors.allowMethods"), ","),
		AllowHeaders:     strings.Join(viper.GetStringSlice("fiber.cors.allowHeaders"), ","),
		AllowCredentials: viper.GetBool("fiber.cors.allowCredentials"),
		ExposeHeaders:    strings.Join(viper.GetStringSlice("fiber.cors.exposeHeaders"), ","),
		MaxAge:           viper.GetInt("fiber.cors.maxAge"),
	}
}
