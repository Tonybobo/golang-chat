package entity

type SignedUrlRequest struct {
	Method string `json:"method" bson:"method"`
	FileName string `json:"fileName" bson:"fileName"`
	ContentType string `json:"contentType" bson:"contentType"`
}