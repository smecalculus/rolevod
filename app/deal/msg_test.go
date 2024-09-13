package deal

import (
	"encoding/json"
	"fmt"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"
)

func TestConverter(t *testing.T) {
	dr := DealRef{
		ID:   id.New[ID](),
		Name: "deal-1",
	}
	x := chnl.Root{
		ID:    id.New[chnl.ID](),
		Name:  "x",
		State: state.OneRef{},
	}
	tran1 := Transition{
		Deal: dr,
		Term: step.Wait{
			X:    chnl.ToRef(x),
			Cont: step.Close{},
		},
	}
	fmt.Printf("core1: %+v\n", tran1)
	mto1 := MsgFromTransition(tran1)
	json1, _ := json.Marshal(mto1)
	fmt.Printf("mto: %v\n", string(json1))
	var mto2 TransitionMsg
	_ = json.Unmarshal(json1, &mto2)
	tran2, _ := MsgToTransition(mto2)
	fmt.Printf("core2: %#v\n", tran2)
}
