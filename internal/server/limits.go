package server
type Tier string
const(TierFree Tier="free";TierPro Tier="pro")
type Limits struct{Tier Tier;Description string}
func LimitsFor(tier string)Limits{if tier=="pro"{return Limits{Tier:TierPro,Description:"Unlimited files, 10GB storage"}}
return Limits{Tier:TierFree,Description:"10 files, 100MB total"}}
func(l Limits)IsPro()bool{return l.Tier==TierPro}
