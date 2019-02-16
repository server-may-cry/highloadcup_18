package accounthelper

import (
	"strconv"
	"strings"
	"time"

	"github.com/server-may-cry/highloadcup_18/structures"
)

var CurrentTime int

func ExtractDomain(email string) string {
	pos := strings.Index(email, "@")
	return email[pos+1:]
}

func IsPremium(p structures.Premium) bool {
	return p.Start < CurrentTime && p.Finish > CurrentTime
}

func WasOrNowPremium(p structures.Premium) bool {
	return p.Start != 0
}

func ExtractPhoneCode(phone string) string {
	if phone == "" {
		return ""
	}
	openPos := strings.Index(phone, "(")
	closePos := strings.Index(phone, ")")
	return phone[openPos+1 : closePos]
}

func ExtractYearOfBirth(ts int) int16 {
	tm := time.Unix(int64(ts), 0)
	return int16(tm.Year())
}

// 1(234)5678901
func PhoneToInt(phone string) int64 {
	if len(phone) == 0 {
		return int64(0)
	}
	//filtered := phone[:1] + phone[2:5] + phone[6:]
	var builder strings.Builder // todo make pool
	builder.WriteByte(phone[0])
	builder.WriteString(phone[2:5])
	builder.WriteString(phone[6:])
	filtered := builder.String()
	i, _ := strconv.ParseInt(filtered, 10, 64)
	return i
}

func PhoneToString(phone int64) string {
	if phone == 0 {
		return ""
	}
	s := strconv.FormatInt(phone, 10)
	return s[:1] + "(" + s[1:4] + ")" + s[4:]
}
