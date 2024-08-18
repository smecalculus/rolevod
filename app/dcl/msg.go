package dcl

type SpecMsg struct {
	Name string `json:"name"`
}

type RootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetMsg struct {
	ID string `param:"id"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgToTpSpec     func(SpecMsg) TpSpec
	MsgFromTpSpec   func(TpSpec) SpecMsg
	MsgToExpSpec    func(SpecMsg) ExpSpec
	MsgFromExpSpec  func(ExpSpec) SpecMsg
	MsgFromTpRoot   func(TpRoot) RootMsg
	MsgFromTpRoots  func([]TpRoot) []RootMsg
	MsgFromExpRoot  func(ExpRoot) RootMsg
	MsgFromExpRoots func([]ExpRoot) []RootMsg
)
