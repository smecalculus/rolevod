package force

type ForceSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type ForceRootMsg struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Children []ForceTeaserMsg `json:"children"`
}

type ForceTeaserMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `json:"name"`
}

type KinshipSpecMsg struct {
	ParentID    string   `param:"id" json:"parent"`
	ChildrenIDs []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   ForceTeaserMsg   `json:"parent"`
	Children []ForceTeaserMsg `json:"children"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	// force
	MsgToForceSpec    func(ForceSpecMsg) (ForceSpec, error)
	MsgFromForceSpec  func(ForceSpec) ForceSpecMsg
	MsgToForceRoot    func(ForceRootMsg) (ForceRoot, error)
	MsgFromForceRoot  func(ForceRoot) ForceRootMsg
	MsgFromForceRoots func([]ForceRoot) []ForceRootMsg
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)
