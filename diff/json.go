package diff

import "encoding/json"

//resourceJson is an intermediate struct for marshalling the given source and the target paths
type resourceJson struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(resourceJson{
		Source: r.source,
		Target: r.target,
	})
}

func (r *Resource) UnmarshalJSON(data []byte) error {
	resJson := resourceJson{}

	err := json.Unmarshal(data, &resJson)
	if err != nil {
		return err
	}

	r.source = resJson.Source
	r.target = resJson.Target

	return nil
}
