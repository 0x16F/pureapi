package wsconnect

// -- get users req --

type methodID int

const (
	MethodIDGetUsers methodID = 9
)

type Gender string

const (
	GenderFemale Gender = "f"
	GenderMale   Gender = "m"
)

type Sexuality string

const (
	SexualityHeterosexual Sexuality = "h"
	SexualityBisexual     Sexuality = "b"
	SexualityGay          Sexuality = "g"
)

type SpokenLanguage string

const (
	SpokenLanguageEnglish SpokenLanguage = "en"
	SpokenLanguageRussian SpokenLanguage = "ru"
)

type RelationshipGoal string

const (
	RelationshipGoalChat RelationshipGoal = "c"
)

type CityID int

const (
	CityIDMoscow CityID = 524901
)

type Radius uint

const (
	Radius5km   Radius = 5
	Radius10km  Radius = 10
	Radius30km  Radius = 30
	Radius50km  Radius = 50
	Radius100km Radius = 100
	RadiusAny   Radius = 0
)

type GetUsersFilters struct {
	Gender            []Gender           `json:"gender,omitempty"`
	Sexuality         []Sexuality        `json:"sexuality,omitempty"`
	SpokenLanguages   []SpokenLanguage   `json:"spoken_languages,omitempty"`
	RelationshipGoals []RelationshipGoal `json:"relationship_goals,omitempty"`
	CityID            CityID             `json:"city_id,omitempty"`
	Radius            Radius             `json:"radius,omitempty"`
}

type smartFeedLogic string

const (
	smartFeedLogic1000kmUsersCountryRadiusNewUsers smartFeedLogic = "1000km-users_country_radius-new_users"
)

type ab struct {
	SmartFeedLogic smartFeedLogic `json:"smart_feed_logic"`
}

type getUsersData struct {
	SessionID string          `json:"session_id"`
	Filters   GetUsersFilters `json:"filters"`
	Ab        ab              `json:"ab"`
}

type paramsMethod string

const (
	paramsMethodSmartFeedRead paramsMethod = "smart_feed.read"
)

type getUsersParams struct {
	Data   getUsersData `json:"data"`
	Method paramsMethod `json:"method"`
}

type getUsersReq struct {
	Method methodID       `json:"method"`
	Params getUsersParams `json:"params"`
	ID     int            `json:"id"`
}

// -- get users resp --

type country struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type region struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Country country `json:"country"`
}

type city struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Region     region `json:"region"`
	NameStd    string `json:"name_std"`
	IsInRadius bool   `json:"is_in_radius"`
}

type reactions struct {
	IncomingLike bool `json:"incoming_like"`
	OutgoingLike bool `json:"outgoing_like"`
	Gift         any  `json:"gift"`
}

type User struct {
	IsRevoked         bool      `json:"is_revoked"`
	UserID            string    `json:"user_id"`
	UserCreatedAt     float64   `json:"user_created_at"`
	Gender            string    `json:"gender"`
	Sexuality         string    `json:"sexuality"`
	Age               any       `json:"age"`
	Height            any       `json:"height"`
	InPair            bool      `json:"in_pair"`
	IsPremiumFeatured bool      `json:"is_premium_featured"`
	DistanceM         int       `json:"distance_m"`
	IsOnline          bool      `json:"is_online"`
	OnlineAt          float64   `json:"online_at"`
	AnnouncementID    string    `json:"announcement_id"`
	AnnouncementText  string    `json:"announcement_text"`
	IsPrefilledText   bool      `json:"is_prefilled_text"`
	Photos            []any     `json:"photos"`
	City              city      `json:"city"`
	Reactions         reactions `json:"reactions"`
	AvatarURL         string    `json:"avatar_url"`
	Temptations       []any     `json:"temptations"`
	SpokenLanguages   []string  `json:"spoken_languages"`
	RelationshipGoal  string    `json:"relationship_goal"`
	HasNewbieBadge    bool      `json:"has_newbie_badge"`
	MatchCondition    string    `json:"match_condition"`
	MlScore           any       `json:"ml_score"`
}

type dataContent struct {
	MatchCondition string `json:"match_condition"`
	IsSuggestions  bool   `json:"is_suggestions"`
	Results        []User `json:"results"`
}

type data struct {
	Success bool        `json:"success"`
	Data    dataContent `json:"data"`
}

type result struct {
	Data data `json:"data"`
}

type getUsersResp struct {
	ID     int    `json:"id"`
	Result result `json:"result"`
}
