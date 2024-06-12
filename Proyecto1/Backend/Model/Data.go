package Model

type DataRAM struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	Used_percentage string `json:"used_percentage"`
	Free_percentage string `json:"free_percentage"`
}

type DataCPU struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	Used_percentage string `json:"used_percentage"`
	Free_percentage string `json:"free_percentage"`
}

type DataProcess struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Pid      string `json:"pid"`
	Name     string `json:"name"`
	User     string `json:"user"`
	State    string `json:"state"`
	Ram      string `json:"ram"`
	PidPadre string `json:"pidPadre"`
}
