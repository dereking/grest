package api

import (
	"errors"
	"time"
)

import (
	"github.com/dereking/grest/security"
)

type ApiAccount struct {
	AppID      string
	AppKey     string
	CreateTime int
	Lifespan   int // days

}

const MAX_TIME_SPAN = 5 //最大时间差:5分钟

var allAPI map[string]*ApiAccount

func FindAPIAccount(appid string) *apiAccount {

	return allAPI[appid]
}

func NewAPIAccount(lifespanDays int) *apiAccount {
	ret := &ApiAccount{
		AppID:      security.GenSessionID(),
		CreateTime: time.Now().UTC().Unix(),
		Lifespan:   lifespanDays * 24 * 3600,
	}

	ret.AppKey = security.GenAppKey(ret.AppID)

}

func (a *ApiAccount) Check(appid, appt, appr, apivc string) error {

	t := time.Now().UTC().Unix()
	at, err := strconv.ParseInt(appt, 10, 64)
	if err != nil {
		return errors.New("ApiAccount.Check : 时间错误")
	}
	if math.Abs(float64(at-t)) > MAX_TIME_SPAN*60 {
		return errors.New(fmt.Sprintf("ApiAccount.Check : 时间差超过%d分钟", MAX_TIME_SPAN))
	}

	appkey := security.GenAppKey(appid)

	vc := security.Md5(fmt.Sprintf("%s_%s_%s_%s", appid, appkey, appt, appr))

	if (len(appid) == 0) || (strings.Compare(vc, apivc) != 0) {

		return errors.New("ApiAccount.Check : 校验错误")
	}
}
