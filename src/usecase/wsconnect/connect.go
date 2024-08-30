package wsconnect

type connectResultData struct {
	SessionID string `json:"client"`
	Version   string `json:"version"`
	Expires   bool   `json:"expires"`
	TTL       int    `json:"ttl"`
}

type connectResp struct {
	ID     int               `json:"id"`
	Result connectResultData `json:"result"`
}
