package zones

type RequestUpdateRRSet struct {
	RRSets []*RRSet `json:"rrsets"`
}
