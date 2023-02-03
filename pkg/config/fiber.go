package config

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
)

type fiberCache struct {
	Expiration time.Duration `default:"1m"`
}
type fiberEncryptCookie struct {
	Key string `default:"secret-thirty-2-character-string"`
}
type fiberLimiter struct {
	Max        int           `default:"5"`
	Expiration time.Duration `default:"1m"`
}
type fiberApp struct {
	Prefork                      bool          `default:"false"`
	ServerHeader                 string        `default:"application" split_words:"true"`
	StrictRouting                bool          `default:"true" split_words:"true"`
	CaseSensitive                bool          `default:"true" split_words:"true"`
	Immutable                    bool          `default:"false"`
	UnescapePath                 bool          `default:"true" split_words:"true"`
	ETag                         bool          `default:"false" envconfig:"ETAG"`
	BodyLimit                    int           `default:"4194304" split_words:"true"`
	Concurrency                  int           `default:"262144"`
	ReadTimeout                  time.Duration `default:"4s" split_words:"true"`
	WriteTimeout                 time.Duration `default:"4s" split_words:"true"`
	IdleTimeout                  time.Duration `default:"20s" split_words:"true"`
	ReadBufferSize               int           `default:"4096" split_words:"true"`
	WriteBufferSize              int           `default:"4096" split_words:"true"`
	CompressedFileSuffix         string        `default:".fiber.gz" split_words:"true"`
	GETOnly                      bool          `default:"false" envconfig:"GET_ONLY"`
	DisableKeepalive             bool          `default:"false" split_words:"true"`
	DisableDefaultDate           bool          `default:"false" split_words:"true"`
	DisableDefaultContentType    bool          `default:"false" split_words:"true"`
	DisableHeaderNormalizing     bool          `default:"false" split_words:"true"`
	DisableStartupMessage        bool          `default:"false" split_words:"true"`
	AppName                      string        `default:"application" split_words:"true"`
	StreamRequestBody            bool          `default:"true" split_words:"true"`
	DisablePreParseMultipartForm bool          `default:"false" split_words:"true"`
	ReduceMemoryUsage            bool          `default:"false" split_words:"true"`
	Network                      string        `default:"tcp4" split_words:"true"`
	EnableTrustedProxyCheck      bool          `default:"true" split_words:"true"`
	TrustedProxies               []string      `default:"[]" split_words:"true"`
	EnableIPValidation           bool          `default:"false" split_words:"true"`
	EnablePrintRoutes            bool          `default:"false" split_words:"true"`
}

var FiberApp fiberApp
var FiberLimiter fiberLimiter
var FiberEncryptCookie fiberEncryptCookie
var FiberCache fiberCache
var Fiber fiber.Config

func LoadFiberConfig() error {
	err := envconfig.Process("fiber", &FiberApp)
	if err != nil {
		return err
	}
	err = envconfig.Process("fiber_limiter", &FiberLimiter)
	if err != nil {
		return err
	}
	err = envconfig.Process("fiber_encrypt_cookie", &FiberEncryptCookie)
	if err != nil {
		return err
	}
	err = envconfig.Process("fiber_cache", &FiberCache)
	if err != nil {
		panic(err)
	}
	Fiber = fiber.Config{
		Prefork:       Fiber.Prefork,
		ServerHeader:  Fiber.ServerHeader,
		StrictRouting: Fiber.StrictRouting,
		CaseSensitive: Fiber.CaseSensitive,
		Immutable:     Fiber.Immutable,
		UnescapePath:  Fiber.UnescapePath,
		ETag:          Fiber.ETag,
		BodyLimit:     Fiber.BodyLimit,
		Concurrency:   Fiber.Concurrency,
		// Views:                nil,
		// ViewsLayout:          "",
		// PassLocalsToViews:    false,
		ReadTimeout:          Fiber.ReadTimeout,
		WriteTimeout:         Fiber.WriteTimeout,
		IdleTimeout:          Fiber.IdleTimeout,
		ReadBufferSize:       Fiber.ReadBufferSize,
		WriteBufferSize:      Fiber.WriteBufferSize,
		CompressedFileSuffix: Fiber.CompressedFileSuffix,
		ProxyHeader:          fiber.HeaderXForwardedFor,
		GETOnly:              Fiber.GETOnly,
		// ErrorHandler: func(*fiber.Ctx, error) error {
		// },
		DisableKeepalive:             Fiber.DisableKeepalive,
		DisableDefaultDate:           Fiber.DisableDefaultDate,
		DisableDefaultContentType:    Fiber.DisableDefaultContentType,
		DisableHeaderNormalizing:     Fiber.DisableHeaderNormalizing,
		DisableStartupMessage:        Fiber.DisableStartupMessage,
		AppName:                      Fiber.AppName,
		StreamRequestBody:            Fiber.StreamRequestBody,
		DisablePreParseMultipartForm: Fiber.DisablePreParseMultipartForm,
		ReduceMemoryUsage:            Fiber.ReduceMemoryUsage,
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		// XMLEncoder: func(v interface{}) ([]byte, error) {
		// },
		Network:                 Fiber.Network,
		EnableTrustedProxyCheck: Fiber.EnableTrustedProxyCheck,
		TrustedProxies:          Fiber.TrustedProxies,
		EnableIPValidation:      Fiber.EnableIPValidation,
		EnablePrintRoutes:       Fiber.EnablePrintRoutes,
		// ColorScheme:             fiber.Colors{},
		// RequestMethods:          []string{},
	}

	return nil
}
