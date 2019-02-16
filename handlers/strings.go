package handlers

//nolint:gochecknoglobals
var (
	emptyBody            = []byte("{}")
	contentType          = []byte("application/json")
	emptyFilterResponse  = []byte(`{"accounts":[]}`)
	emptyGroupBody       = []byte(`{"groups":[]}`)
	emptyRecommendBody   = []byte(`{"accounts":[]}`)
	emptySuggestResponse = []byte(`{"accounts":[]}`)
)
