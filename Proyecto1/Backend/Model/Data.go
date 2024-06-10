package Model

type Data struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	Used_percentage string `json:"used_percentage"`
	Free_percentage string `json:"free_percentage"`
}
