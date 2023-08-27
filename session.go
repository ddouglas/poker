package poker

type Session struct {
	ID     string
	UserID *string `dynamodbav:",omitempty"`
	State  *string `dynamodbav:",omitempty"`
}
