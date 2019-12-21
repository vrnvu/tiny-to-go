package tiny

type Redirect struct {
	Code      string `json: "code" bson:"code" msgpack: "code"`
	URL       string `json: "url" bson:"url" msgpack: "url" validate:"empty=false & format=url"`
	Timestamp int64  `json: "timestamp" bson:"timestamp" msgpack: "timestamp"`
}
