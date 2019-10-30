package util

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06",   //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01",      // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1",       // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan",     // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2",  // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon",    // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3",  // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

///////////////////////////////////////////////////////////////////////////////
// something related to try catch finally
// demo:
// Try(panic("123"))
// .Catch(1, func(e interface{}) { fmt.Println("int", e)})
// .Catch("", func(e interface{}) {fmt.Println("string", e)})
// .Finally(func() {fmt.Println("error")})
type TryExceptionHandler func(interface{})

type tryStruct struct {
	catches map[reflect.Type]TryExceptionHandler
	hold    func()
}

func Try (f func()) *tryStruct {
	return &tryStruct {
		catches: make(map[reflect.Type]TryExceptionHandler),
		hold:    f,
	}
}

// register the exception func (second parameter) with the reflect type of first parameter
func (t *tryStruct) Catch(e interface{}, f TryExceptionHandler) *tryStruct {
	t.catches[reflect.TypeOf(e)] = f
	return t
}

//do the call and if something go wrong, it will try to recover and exec the func which was setted at the catch stage
func (t *tryStruct) Finally(f func()) {
	defer func() {
		if e:= recover(); nil != e {
			if h, ok := t.catches[reflect.TypeOf(e)]; ok {
				h(e)
			} else {
				f()
			}
		}
	}()
	t.hold()
}

///////////////////////////////////////////////////////////////////////////////

func Assert(err error) {
	if nil != err {
		panic(err)
	}
}

func Debug() {

}

func Int2str(num int) (string) {
	return strconv.Itoa(num)
}

func Str2int(str string) (int) {
	num, err := strconv.Atoi(str)
	if nil != err {
		num = 0
	}

	return  num
}

func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

/**
 * @param str, string number.
 * @param defaultNum, return this when str conv into in err.
 * @return int
 */
func Str2int2(str string, defaultNum int) (int) {
	num, err := strconv.Atoi(str)

	if nil != err {
		num = defaultNum
	}

	return num
}

func Int642str(num int64) (string) {
	return strconv.FormatInt(num,10)
}

func Str2int64 (str string) (int64) {
	num, err := strconv.ParseInt(str, 10, 64)

	if nil != err {
		num = 0
	}

	return num
}

func Str2int642 (str string, defaultNum int64) (int64) {
	num, err := strconv.ParseInt(str, 10, 64)

	if nil != err {
		num = defaultNum
	}

	return num
}

func DateFormat(t time.Time, format string) string {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return t.Format(format)
}

func Date2unix(date string) int64 {
	timezone, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", date, timezone)

	return tmp.Unix()
}

func Unix2time(utime int64) string {
	return time.Unix(utime, 0).Format("2006-01-02 15:04:05")
}

func Unix2date(utime int64) string {
	return time.Unix(utime, 0).Format("2006-01-02")
}

func Unix2year(utine int64) string {
	return time.Unix(utine, 0).Format("2006")
}

func Str2byte(str string) []byte {
	return []byte(str)
}

func Byte2str(b []byte) string {
	return string(b)
}

func Byte2int(b []byte)int{
	var ret int = 0
	var len int = len(b)
	var i uint = 0
	for i=0; i<uint(len); i++{
		ret = ret | (int(b[i]) << (i*8))
	}
	return ret
}

func Int2byte(i int) []byte {
	var len uintptr = unsafe.Sizeof(i)
	ret := make([]byte, len)
	var tmp int = 0xff
	var index uint = 0
	for index=0; index < uint(len); index++ {
		ret[index] = byte((tmp << (index*8) & i) >> (index*8))
	}

	return ret
}

func OnceTimerTask(second time.Duration, f func()) {
	timer := time.NewTimer(time.Second * second)
	go func() {
		<- timer.C
		f()
	}()
}

func Join(arr []string, flag string) string {
	ret := ""
	if len(arr) < 1 {
		return ret
	}

	for _, v := range arr {
		ret += v + flag
	}

	return ret[0:len(ret) - 1]
}

func Request(reqType string, url string, params string, headers map[string]string) (string,error) {
	ret := ""
	client := http.Client{}
	req, err := http.NewRequest(strings.ToUpper(reqType), url, strings.NewReader(params))
	if nil != err {
		return ret, err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if nil != err {
		return ret, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return ret, err
	}

	ret = string(body)

	return ret,nil
}