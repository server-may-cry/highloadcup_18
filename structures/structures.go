package structures

const (
	StatusFree       = "свободны"
	StatusHold       = "заняты"
	StatusAllComplex = "всё сложно"
	SexF             = "f"
	SexM             = "m"
)

const (
	StatusFreeInt int8 = iota
	StatusHoldInt
	StatusAllComplexInt
	SexFInt int8 = iota
	SexMInt
)

type Account struct {
	ID        int      `json:"id"`
	Fname     string   `json:"fname"`
	Sname     string   `json:"sname"`
	Email     string   `json:"email"`
	Interests []string `json:"interests"`
	Status    string   `json:"status"`
	Premium   Premium  `json:"premium"`
	Sex       string   `json:"sex"`
	Phone     string   `json:"phone"` // uniq
	Likes     []Like   `json:"likes"`
	Birth     int      `json:"birth"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Joined    int      `json:"joined"`
}

type Premium struct {
	Start  int `json:"start"`
	Finish int `json:"finish"`
}

type Like struct {
	Ts int `json:"ts"`
	ID int `json:"id"`
}
