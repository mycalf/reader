package converter

import (
	"regexp"
	"strings"
	"unicode"

	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/net/html/charset"
)

// Content Struct
// 内容结构 ...
type Content struct {
	HTML  string
	Bytes []byte
}

/*---------------------------------------------------------------*/

// Run Converter Function
// 运行对当前进程进行编码转换成UTF-8 ...
func (content *Content) Run() (string, bool) {

	// 自动获取资源编码 ...
	charset, ok := content.getCharset()

	// 未获取到资源编码 ...
	if !ok {
		return "", false
	}

	// UTF-8无需转换 ...
	if charset == "UTF-8" {
		return content.HTML, true
	}

	// 转换其他编码至UTF-8 ...
	if cil, err := iconv.NewConverter(charset, "UTF-8"); err == nil {

		defer cil.Close()

		// 开始转换编码 ...
		output, err := cil.ConvertString(content.HTML)

		if err == nil {
			return output, true
		}
	}

	return "", false
}

/*---------------------------------------------------------------*/

// Charset Function
// 返回当前进程的字符集 ...
func (content *Content) getCharset() (string, bool) {

	// 自动获取编码 ...
	encoding, name, ok := charset.DetermineEncoding(content.Bytes, "")

	// 如果自动获取成功或encoding不为空
	// 则输出编码格式 ...
	if ok {
		return strings.ToUpper(name), true
	}

	if encoding != nil && name != "windows-1252" {
		return strings.ToUpper(name), true
	}

	// 如果内容中出现汉字
	// 则输出GB18030 ...
	if isHan(content.HTML) {
		return "GB18030", true
	}

	// 不符合上述条件
	// 则返回空 ...
	return "", false
}

/*---------------------------------------------------------------*/

// IsHan Function
// 判断是否存在中文 ...
func isHan(str string) bool {
	hanLen := len(regexp.MustCompile("[\\P{Han}]").ReplaceAllString(str, ""))
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || hanLen > 0 {
			return true
		}
	}
	return false
}

/*---------------------------------------------------------------*/

/* End of file converter.go */
/* Location: ./converter.go */
