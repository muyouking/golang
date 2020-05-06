package Tools

//其它工具类
//md5,中文转码
import (
	"bufio"
	"crypto/md5"
	// "hash/crc32"
	"fmt"
	"io"
	"os"

	"github.com/axgle/mahonia"
)

//生成并返回Md5码
func Md5sum3(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//处理中文字符串的方法
func EncoderGBK(str string) string {

	return mahonia.NewEncoder("gbk").ConvertString(str)

}

//mahonia.NewEncoder 编码
func EncoderUTF8(str string) string {
	return mahonia.NewEncoder("UTF-8").ConvertString(str)
}

//mahonia.NewDecoder 解码
func DecoderGBK(str string) string {
	return mahonia.NewDecoder("gbk").ConvertString(str)

}

func DecoderUTF8(str string) string {
	return mahonia.NewDecoder("UTF-8").ConvertString(str)

}
