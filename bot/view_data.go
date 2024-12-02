package bot

import (
	"time"

	macaron "gopkg.in/macaron.v1"
)

func viewData(ctx *macaron.Context) {
	var r *User
	dr := &DataResponse{}
	tgid := getTgId(ctx)
	ref := ctx.Params("referral")
	code := ctx.Params("code")
	name := ctx.Params("name")

	if tgid != 0 {
		u := getUserOrCreate2(tgid, code, name)
		r = getUserByCode(ref)

		if u.ReferrerID == nil && r.ID != u.ID && r.ID != 0 {
			u.ReferrerID = &r.ID
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
			notify(lNewRef, r.TelegramId)
		}

		dr.Code = u.Code
		dr.AddressDeposit = u.AddressDeposit
		dr.AddressWithdraw = u.AddressWithdraw
		dr.TMU = float64(u.TMU) / float64(Mul9)
		dr.Earnings = float64(u.rewards()) / float64(Mul9)
		dr.LastUpdated = u.LastUpdated
		dr.TimeLock = u.TimeLock
		dr.IsFollower = u.isFollower()
		dr.IsMember = u.isMember()
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, dr)
}

type DataResponse struct {
	Earnings        float64    `json:"earnings"`
	TMU             float64    `json:"tmu"`
	Code            string     `json:"code"`
	AddressDeposit  string     `json:"addr_deposit"`
	AddressWithdraw string     `json:"addr_withdraw"`
	LastUpdated     time.Time  `json:"last_updated"`
	TimeLock        *time.Time `json:"time_lock"`
	IsFollower      bool       `json:"is_follower"`
	IsMember        bool       `json:"is_member"`
}
