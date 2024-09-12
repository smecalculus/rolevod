package chnl

type RefData struct {
	ID   string `db:"id" json:"id,omitempty" `
	Name string `db:"name" json:"name,omitempty"`
}

type rootData struct {
	ID    string `db:"id"`
	PreID string `db:"pre_id"`
	Name  string `db:"name"`
	State string `db:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	DataToRef    func(*RefData) (Ref, error)
	DataFromRef  func(Ref) *RefData
	DataToRefs   func([]RefData) ([]Ref, error)
	DataFromRefs func([]Ref) []RefData
	DataToRoot   func(rootData) (Root, error)
	DataFromRoot func(Root) rootData
)
