package prog

import (
	"fmt"

	sdkpooltypedef "github.com/orglang/go-sdk/pool/typedef"
	sdkproctypedef "github.com/orglang/go-sdk/proc/typedef"
	"github.com/orglang/go-sdk/prog"

	pooltypedef "orglang/go-engine/pool/typedef"
	proctypedef "orglang/go-engine/proc/typedef"
)

func MsgToSpec(dto prog.Spec) (Spec, error) {
	spec := Spec{}
	for _, msgExp := range dto.Exps {
		switch msgExp := msgExp.(type) {
		case sdkpooltypedef.DefSpec:
			poolDef, convErr := pooltypedef.MsgToDefSpec(msgExp)
			if convErr != nil {
				return Spec{}, convErr
			}
			spec.Pools = append(spec.Pools, poolDef)
		case sdkproctypedef.DefSpec:
			procDef, convErr := proctypedef.MsgToDefSpec(msgExp)
			if convErr != nil {
				return Spec{}, convErr
			}
			spec.Procs = append(spec.Procs, procDef)
		default:
			return Spec{}, fmt.Errorf("unsupported program expression: %T", msgExp)
		}
	}
	return spec, nil
}
