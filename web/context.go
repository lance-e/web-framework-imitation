package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req *http.Request

	// Resp 如果用户直接使用这个
	//那么他们将绕开 RespData 和 RespStatusCode 两个
	//部分 middleware 将无法运作
	Resp http.ResponseWriter

	//主要是给 middleware 读写用的
	RespData       []byte
	RespStatusCode int

	PathParams map[string]string

	queryValues url.Values //缓存的数据

	MatchedRoute string //命中的路由
}

type StringValue struct {
	value string
	err   error
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}
func (c *Context) BindJSON(value any) error {
	if value == nil {
		return errors.New("web :输入不能为nil")
	}
	if c.Req.Body == nil {
		return errors.New("web :body不能为nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(value)
}

// Form :URL里面的查询参数和PATCH，POST，PUT的表单数据。 所有表单数据都可以拿到
// PostForm：PATCH，POST，PUT body参数，在编码是x-www-form-urlencoded的时候才能拿到
func (c *Context) FormValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{err: err}
	}
	return StringValue{value: c.Req.FormValue(key)}
}

// Query 和form的区别就是没有缓存
func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	values, ok := c.Req.URL.Query()[key]
	if !ok {
		return StringValue{err: errors.New("web :找不到这个key")}
	}
	return StringValue{value: values[0]}
	//用户区分不出来是真的有值，但是当值为一个空字符串，就不知道是不是真的有值
	//return c.queryValues.Get(key), nil
}
func (c *Context) PathValue(key string) StringValue {
	value, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web :key 不存在")}
	}
	return StringValue{value: value}
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.value, 10, 64)
}
func (c *Context) RespJSON(status int, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.RespData = data
	c.RespStatusCode = status
	return err
}
