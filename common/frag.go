package common

import (
	"encoding/json"
	"fmt"
)

//Frag represent detail info about a frag
type Frag struct {
	Index uint64 `json:"index" bson:"index,omitempty"`
	Start int64  `json:"start" bson:"start,omitempty"`
	Size  int64  `json:"size" bson:"size,omitempty"`
}

func (f *Frag) String() string {
	return fmt.Sprintf("Frag - index: %d, start: %d, size: %d", f.Index, f.Start, f.Size)
}

func MustUnmarshalJsonToFrag(b []byte) Frag {
	var frag Frag
	if err := json.Unmarshal(b, &frag); err != nil {
		panic(err)
	}
	return frag
}

func MustMarshalFragAsJson(f Frag) []byte {
	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	return b
}
