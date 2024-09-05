package work

type WorkSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type WorkRootMsg struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Children []WorkTeaserMsg `json:"children"`
}

type WorkTeaserMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `json:"name"`
}

type KinshipSpecMsg struct {
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   WorkTeaserMsg   `json:"parent"`
	Children []WorkTeaserMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// work
	MsgToWorkSpec    func(WorkSpecMsg) (WorkSpec, error)
	MsgFromWorkSpec  func(WorkSpec) WorkSpecMsg
	MsgToWorkRoot    func(WorkRootMsg) (WorkRoot, error)
	MsgFromWorkRoot  func(WorkRoot) WorkRootMsg
	MsgFromWorkRoots func([]WorkRoot) []WorkRootMsg
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
