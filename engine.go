package goecharts

import (
	"bytes"
	"github.com/gobuffalo/packr"
	"html/template"
	"io"
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	LETTER_BYTES    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LETTER_IDX_BITS = 6
	LETTER_IDX_MASK = 1<<LETTER_IDX_BITS - 1 // All 1-bits, as many as letterIdxBits
	LETTER_IDX_MAX  = 63 / LETTER_IDX_BITS
	CHART_SIZE      = 12
)

// TODO: template must
// 渲染图表
func renderChart(chart interface{}, w io.Writer) {
	box := packr.NewBox("./templates")
	htmlContent, err := box.FindString("index.html")
	t, err := template.New("").Parse(htmlContent)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, chart)
}

var seed = rand.NewSource(time.Now().UnixNano())

// 生成唯一且随机的图表 ID
func genChartID() string {
	b := make([]byte, CHART_SIZE)
	for i, cache, remain := CHART_SIZE-1, seed.Int63(), LETTER_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = seed.Int63(), LETTER_IDX_MAX
		}
		if idx := int(cache & LETTER_IDX_MASK); idx < len(LETTER_BYTES) {
			b[i] = LETTER_BYTES[idx]
			i--
		}
		cache >>= LETTER_IDX_BITS
		remain--
	}
	return string(b)
}

// 过滤替换渲染结果
func replaceRender(b bytes.Buffer) []byte {
	// __x__ 与模板占位符相匹配
	pat, err := regexp.Compile(`(__x__")|("__x__)`)
	if err != nil {
		log.Println(err)
	}
	// 替换并转为 []byte 类型
	res := []byte(pat.ReplaceAllString(b.String(), "_x_"))
	return res
}

// 为结构体设置默认值
// 部分代码参考 https://github.com/mcuadros/go-defaults
func setDefaultValue(ptr interface{}) error {
	var err error
	//需要参数为指针类型，指针才能改变值
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return err
	}

	elem := reflect.ValueOf(ptr).Elem()
	t := elem.Type()

	//判断是否为结构体
	if t.Kind() != reflect.Struct {
		return err
	}

	for i := 0; i < t.NumField(); i++ {
		//如果没有 `default` tag 则不作处理
		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "" {
			setField(elem.Field(i), defaultVal)
		}
	}
	return nil
}

// 为具体字段设置默认值
func setField(field reflect.Value, defaultVal string) {
	// 目前只判断 string, bool 两种变量类型
	switch field.Kind() {
	// string 类型
	case reflect.String:
		if field.String() == "" {
			field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
		}
		// bool 类型
	case reflect.Bool:
		if val, err := strconv.ParseBool(defaultVal); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	}
}
