package fiberlogger

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// FuncTag is a function used to populate logrus field
type FuncTag func(c *fiber.Ctx, d *data) interface{}

// predefined FuncTag functions
var (
	FuncTagReferer FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Get(fiber.HeaderReferer)
	}
	FuncTagProtocol FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Protocol()
	}
	FuncTagPort FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Port()
	}
	FuncTagIP FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.IP()
	}
	FuncTagIPs FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Get(fiber.HeaderXForwardedFor)
	}
	FuncTagHost FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Hostname()
	}
	FuncTagPath FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Path()
	}
	FuncTagURL FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.OriginalURL()
	}
	FuncTagUA FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Get(fiber.HeaderUserAgent)
	}
	FuncTagBody FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Body()
	}
	FuncTagBytesReceived FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return len(c.Request().Body())
	}
	FuncTagBytesSent FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		if c.Response().Header.ContentLength() < 0 {
			return 0
		}
		return len(c.Response().Body())
	}
	FuncTagRoute FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Route().Path
	}
	FuncTagResBody FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Response().Body()
	}
	FuncTagReqHeaders FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		reqHeaders := make([]string, 0)
		for k, v := range c.GetReqHeaders() {
			reqHeaders = append(reqHeaders, k+"="+v)
		}
		return []byte(strings.Join(reqHeaders, "&"))
	}
	FuncTagQueryStringParams FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Request().URI().QueryArgs().String()
	}
	FuncTagStatus FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Response().StatusCode()
	}
	FuncTagMethod FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return c.Method()
	}
	FuncTagPid FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return d.pid
	}
	FuncTagLatency FuncTag = func(c *fiber.Ctx, d *data) interface{} {
		return d.end.Sub(d.start).String()
	}
	FuncTagReqHeader = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			return c.Get(extra)
		}
	}
	FuncTagRespHeader = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			return c.GetRespHeader(extra)
		}
	}
	FuncTagQuery = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			return c.Query(extra)
		}
	}
	FuncTagForm = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			return c.FormValue(extra)
		}
	}
	FuncTagCookie = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			return c.Cookies(extra)
		}
	}
	FuncTagLocals = func(extra string) FuncTag {
		return func(c *fiber.Ctx, d *data) interface{} {
			switch v := c.Locals(extra).(type) {
			case []byte:
				return string(v)
			case string:
				return v
			case nil:
				return nil
			default:
				return fmt.Sprintf("%v", v)
			}
		}
	}
)

// attached keyTag separator
const sep string = ":"

// getFuncTagMap selects functions to be used for logrus fields population
func getFuncTagMap(cfg Config, d *data) map[string]FuncTag {
	m := make(map[string]FuncTag, len(cfg.Tags))
	for _, t := range cfg.Tags {
		switch t {
		case TagReferer:
			m[TagReferer] = FuncTagReferer
		case TagProtocol:
			m[TagProtocol] = FuncTagProtocol
		case TagPort:
			m[TagPort] = FuncTagPort
		case TagIP:
			m[TagIP] = FuncTagIP
		case TagIPs:
			m[TagIPs] = FuncTagIPs
		case TagHost:
			m[TagHost] = FuncTagHost
		case TagPath:
			m[TagPath] = FuncTagPath
		case TagURL:
			m[TagURL] = FuncTagURL
		case TagUA:
			m[TagUA] = FuncTagUA
		case TagBody:
			m[TagBody] = FuncTagBody
		case TagBytesReceived:
			m[TagBytesReceived] = FuncTagBytesReceived
		case TagBytesSent:
			m[TagBytesSent] = FuncTagBytesSent
		case TagRoute:
			m[TagRoute] = FuncTagRoute
		case TagResBody:
			m[TagResBody] = FuncTagResBody
		case TagReqHeaders:
			m[TagReqHeaders] = FuncTagReqHeaders
		case TagQueryStringParams:
			m[TagQueryStringParams] = FuncTagQueryStringParams
		case TagStatus:
			m[TagStatus] = FuncTagStatus
		case TagMethod:
			m[TagMethod] = FuncTagMethod
		case TagPid:
			m[TagPid] = FuncTagPid
		case TagLatency:
			m[TagLatency] = FuncTagLatency
		default:
			for _, v := range KeyTags {
				if strings.Contains(t, v) {
					a := strings.Split(t, sep)
					switch a[0] {
					case TagReqHeader:
						m[TagReqHeader] = FuncTagReqHeader(a[1])
					case TagRespHeader:
						m[TagRespHeader] = FuncTagRespHeader(a[1])
					case TagQuery:
						m[TagQuery] = FuncTagQuery(a[1])
					case TagForm:
						m[TagForm] = FuncTagForm(a[1])
					case TagCookie:
						m[TagCookie] = FuncTagCookie(a[1])
					case TagLocals:
						m[TagLocals] = FuncTagLocals(a[1])
					}

				}
			}
		}
	}
	return m
}
