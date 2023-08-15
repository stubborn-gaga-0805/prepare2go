package helpers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/jaevor/go-nanoid"
	"github.com/mozillazg/go-pinyin"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/consts"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ToInt8(i interface{}) int8 {
	return cast.ToInt8(i)
}

func ToInt16(i interface{}) int16 {
	return cast.ToInt16(i)
}

func ToInt32(i interface{}) int32 {
	return cast.ToInt32(i)
}

func ToInt64(i interface{}) int64 {
	return cast.ToInt64(i)
}

func IntToString(i interface{}) string {
	return cast.ToString(i)
}

func StringToInt(i interface{}) int {
	return cast.ToInt(i)
}

func Int64ToStringSlice(slice []int64) []string {
	res := make([]string, 0)
	for _, v := range slice {
		res = append(res, IntToString(v))
	}
	return res
}

func StringToInt64Slice(slice []string) []int64 {
	res := make([]int64, 0)
	for _, v := range slice {
		res = append(res, ToInt64(v))
	}
	return res
}

// CheckPasswordMd5 md5方式检查密码
func CheckPasswordMd5(password string, encrypt string) bool {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil)) == encrypt
}

// Md5Encrypt 用MD5加密字符串
func Md5Encrypt(src string) (encrypt string) {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

// Yuan2Cent 元转分
func Yuan2Cent(yuan string) int32 {
	float, err := strconv.ParseFloat(yuan, 64)
	if err != nil {
		return 0
	}
	return cast.ToInt32(float * 100)
}

// Cent2Yuan 分转元
func Cent2Yuan(cent int32) string {
	return fmt.Sprintf("%.2f", float32(cent)/float32(100))
}

// GenSerialNo 生成序列号
func GenSerialNo(len int) string {
	return fmt.Sprintf("%s%s", time.Now().Format("200601"), nanoId(consts.NumericAlphabet, len-6))
}

// GenUserNo 生成用户编号
func GenUserNo(len int) string {
	return nanoId(consts.UserNoAlphabet, len)
}

// GenNanoId 基于NanoId生成序列号
func GenNanoId(len int) string {
	return nanoId(consts.NormalAlphabet, len)
}

func nanoId(alpha string, len int) string {
	nanoId, err := nanoid.CustomUnicode(alpha, len)
	if err != nil {
		return ""
	}

	return nanoId()
}

// GenUUID 生成UUID
func GenUUID() string {
	return uuid.NewV4().String()
}

// SplitWords 分词
/*func SplitWords(sentence string) (words []string) {
	// 去掉转义、标点符号和emoji
	dict := []string{
		//"./dict/hmm_model.utf8",
		//"./dict/idf.utf8",
		"./dict/jieba.dict.utf8",
		//"./dict/stop_words.utf8",
		"./dict/user.dict.utf8",
	}
	sentence = RemovePunctuation(RemoveEmoji(regexp.MustCompile(`\\.`).ReplaceAllString(sentence, "")))
	sw := gojieba.NewJieba(dict...)
	//defer sw.Free()

	// 分词
	return sw.Cut(sentence, true)
}*/

// RemoveEmoji 去掉字符串的Emoji
func RemoveEmoji(sentence string) string {
	var result []rune
	for _, r := range sentence {
		if !unicode.IsGraphic(r) {
			continue
		}
		if unicode.Is(unicode.So, r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

// RemovePunctuation 去掉字符串的标点符号
func RemovePunctuation(sentence string) string {
	sentence = strings.TrimFunc(sentence, func(r rune) bool {
		return unicode.IsPunct(r)
	})

	var result []rune
	for _, r := range sentence {
		if unicode.IsPunct(r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

// RemoveDuplicates 去掉切片中的重复内容
func RemoveDuplicates(s []string) []string {
	m := make(map[string]bool)
	for _, v := range s {
		m[v] = true
	}
	var result []string
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Han2Pinyin 汉字转拼音
func Han2Pinyin(han string) string {
	pyList := make([]string, 0)
	for _, r := range pinyin.Pinyin(RemovePunctuation(RemoveEmoji(han)), pinyin.NewArgs()) {
		pyList = append(pyList, r[0])
	}
	return strings.Join(pyList, "")
}

// RandomInt 产生随机数
func RandomInt(start, end int) int {
	rand.Seed(time.Now().UnixNano())
	return start + rand.Intn(end-start+1)
}

// Base64Img2Byte Base64图片编码转[]byte
func Base64Img2Byte(base64Img string) (byte []byte, err error) {
	var (
		strList    = strings.Split(base64Img, ",")
		encodedStr = strList[0]
	)
	if len(strList) == 2 {
		encodedStr = strList[1]
	}

	// 将base64字符串解码为[]byte
	decoded, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		zap.S().Errorf("base64.StdEncoding.DecodeString decoded: %v Err: %v", decoded, err)
		return nil, err
	}
	//// 将[]byte解码为png图片
	//img, _, err := image.Decode(strings.NewReader(string(decoded)))
	//if err != nil {
	//	zap.S().Errorf("image.Decode(strings.NewReader(string(decoded))) Err: %v", err)
	//	return nil, err
	//}
	//// 将png图片转换为[]byte
	//buf := new(bytes.Buffer)
	//err = png.Encode(buf, img)
	//if err != nil {
	//	zap.S().Errorf("png.Encode(buf, img) Err: %v", err)
	//	return nil, err
	//}
	//
	//return buf.Bytes(), nil
	return decoded, nil
}

func IsValidPhoneNum(mobile string) bool {
	reg := `^1[3-9]\d{9}$`
	pattern := regexp.MustCompile(reg)
	return pattern.MatchString(mobile)
}

// ParseStartAndEnd 获取开始时间的 00:00:00 和结束时间的 23:59:59
func ParseStartAndEnd(start time.Time, end time.Time) (startAt time.Time, endAt time.Time) {
	startAt = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	if end.IsZero() {
		endAt = startAt.AddDate(0, 0, 1).Add(-time.Second)
	} else {
		endAt = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	}
	return
}

func ParseStartAndEnd2Str(start time.Time, end time.Time) (string, string) {
	start, end = ParseStartAndEnd(start, end)
	return start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05")
}

func NewProgressBar(count int, msg string) (bar *pb.ProgressBar) {
	var progressTemp = `{{string . "prefix" | blue}} {{ bar . "[" "=" (cycle . "↖" "↗" "↘" "↙" ">" ">" ">") "-" "]"}} {{percent . | blue}} {{speed . | blue }}   {{string . "duration" | green}} {{etime . | green}} {{string . "end"}}`

	bar = pb.New(count)
	bar.SetTemplate(pb.ProgressBarTemplate(progressTemp)).
		Set("prefix", msg).
		Set("end", "\n").
		Set("duration", "耗时:").
		SetRefreshRate(time.Second * 10).
		SetWidth(160).
		SetWriter(os.Stdout).
		Start()

	return bar
}

func ContextWithRequestId(ctx context.Context) context.Context {
	return context.WithValue(ctx, consts.ContextRequestIdKey, GenUUID())
}

func GetContextWithRequestId() context.Context {
	return context.WithValue(context.Background(), consts.ContextRequestIdKey, GenUUID())
}

func GetRequestIdFromContext(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(consts.ContextRequestIdKey))
}

func DivCeil(a, b int) int {
	return int(math.Ceil(float64(a) / float64(b)))
}

func WsJsonMsg2Map(jsonMsg interface{}) (res map[string]interface{}, err error) {
	jsonByte, err := json.Marshal(jsonMsg)
	if err != nil {
		return nil, err
	}
	res = make(map[string]interface{}, 0)
	if err = json.Unmarshal(jsonByte, &res); err != nil {
		return nil, err
	}
	return res, nil
}
