// This package hold enum values or custom type of wallet service response code and status
package contract

import (
	"fmt"
	"sort"
	"strings"
)

type (
	// StatusCode custom type to hold value for standard status code.
	StatusCode uint
)

// Constant standar int.
const ssc = 2100

// Status Code & Message enum.
const (
	OK StatusCode = (iota + ssc)
	INSUFFICIENT
	TRXEXISTS
	CANCELEDTRXNOTEXISTS
	INTERNALERROR
	VALIDATIONERROR
	CTXVALNOTFOUNDMETADATANULL
	WALLETNOTFOUND
	QUEUETASKEXISTS
	WALLETALREADYEXISTS
	WALLETLOCKED
	WALLETDISABLED
	WALLETPROMONOTFOUND
	METHODNOTSUPPORTED
	TOOMANYREQUESTS
	ERRORUNLOCKMUTEXREDIS

	enumStatusLimit
)

// String give StatusCode type string value
func (s StatusCode) String(custom ...string) string {
	if len(custom) > 0 {
		return custom[0]
	}

	return [...]string{
		"OK",
		"Insufficient Balance",
		"TransactionId Already Exists",
		"Canceled Transaction Not Exists",
		"Internal Error",
		"Validation Error {0}",
		"Metadata is NULL or Context Not Found",
		"Wallet Not Found",
		"Queue Task Already Exists",
		"Wallet Already Exists",
		"Wallet Locked",
		"Wallet Disabled",
		"Wallet Promo Not Found",
		"Method Not Supported",
		"Too Many Requests",
		"Error Unlock Mutex Redis",
	}[s-ssc]
}

// String give StatusCode type string value
func (s StatusCode) FormatedString(appended ...string) (status string) {
	for k, v := range appended {
		status = strings.ReplaceAll(s.String(), fmt.Sprintf("{%d}", k), v)
	}

	return
}

// StatusCodeLists return all status code enum into map[uint]string list.
func StatusCodeLists() (l map[uint]string) {
	l = make(map[uint]string, enumStatusLimit)
	for i := StatusCode(0 + ssc); i < enumStatusLimit; i++ {
		l[uint(i)] = i.String()
	}

	return
}

func StatusCodeToMDTable() (mdt string) {
	mdTbHeader := `| Code | Status | Desc |`
	l := StatusCodeLists()
	keys := make([]int, 0)
	for k := range l {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	mdt = mdTbHeader
	mdt += fmt.Sprintln()
	mdt += `:---: | --- | ---`
	for _, k := range keys {
		mdt += fmt.Sprintln()

		var d string
		switch k {
		case int(TRXEXISTS):
			d = "gRPC will return success response, but response code will be appear this."

		case int(INTERNALERROR):
			d = "Sometimes will followed with error message. i.e: Deadlock context, etc."

		case int(VALIDATIONERROR):
			d = "{} Will be replaced with the validation message. i.e: payload amount must be absolute value."

		case int(QUEUETASKEXISTS):
			d = "This response status code only appear in 'Internal System' call not 'Provider System' call. i.e: Depo/WD async, Referral release, Rebate release, etc."

		case int(WALLETLOCKED):
			// member can't wd, but still can check balance and wager / bonus provider
			d = "Once wallet is locked, member cannot withdraw, but still can check balance and place a wager"

		case int(WALLETDISABLED):
			d = "Once wallet is disabled, member cannot withdraw, place a wager and check balance"
		}

		mdt += fmt.Sprintf(`| %v | %s | %s |`, k, l[uint(k)], d)
	}

	return
}
