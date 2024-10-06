package deal_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/deal"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/seat"
)

var (
	roleApi = role.NewRoleApi()
	seatApi = seat.NewSeatApi()
	dealApi = deal.NewDealApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestEstablishKinship(t *testing.T) {
	// given
	ps := deal.DealSpec{Name: "parent-deal"}
	pr, err := dealApi.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := deal.DealSpec{Name: "child-deal"}
	cr, err := dealApi.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := deal.KinshipSpec{
		ParentID: pr.ID,
		ChildIDs: []id.ADT{cr.ID},
	}
	err = dealApi.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := dealApi.Retrieve(pr.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedChild := deal.ConvertRootToRef(cr)
	if !slices.Contains(actual.Children, expectedChild) {
		t.Errorf("unexpected children in %q; want: %+v, got: %+v", pr.Name, expectedChild, actual.Children)
	}
}

func TestTakeWaitClose(t *testing.T) {
	// given
	roleSpec := role.RoleSpec{
		FQN: "role-1",
		St:  state.OneSpec{},
	}
	roleRoot, err := roleApi.Create(roleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: roleRoot.St.RID(),
			St:   roleRoot.St,
		},
	}
	seatRoot1, err := seatApi.Create(seatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Ctx: []chnl.Spec{
			seatSpec1.Via,
		},
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: roleRoot.St.RID(),
			St:   roleRoot.St,
		},
	}
	seatRoot2, err := seatApi.Create(seatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	providerProc, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: providerProc.PID,
		},
	}
	clientProc, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerProcPID := clientProc.Ctx[seatSpec1.Via.Name]
	// and
	closeSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: providerProc.PID,
		Term: step.CloseSpec{
			A: providerProcPID,
		},
	}
	// when
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: clientProc.PID,
		Term: step.WaitSpec{
			X: providerProcPID,
			Cont: step.CloseSpec{
				A: clientProc.PID,
			},
		},
	}
	// and
	err = dealApi.Take(waitSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func TestTakeRecvSend(t *testing.T) {
	// given
	roleSpec1 := role.RoleSpec{
		FQN: "role-1",
		St: state.LolliSpec{
			Y: state.OneSpec{},
			Z: state.OneSpec{},
		},
	}
	roleRoot1, err := roleApi.Create(roleSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	roleSpec2 := role.RoleSpec{
		FQN: "role-2",
		St:  state.OneSpec{},
	}
	roleRoot2, err := roleApi.Create(roleSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: roleRoot1.St.RID(),
			St:   roleRoot1.St,
		},
	}
	seatRoot1, err := seatApi.Create(seatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: roleRoot2.St.RID(),
			St:   roleRoot2.St,
		},
		Ctx: []chnl.Spec{
			seatSpec1.Via,
		},
	}
	seatRoot2, err := seatApi.Create(seatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	providerProc, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: providerProc.PID,
		},
	}
	clientProc, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerProcPID := clientProc.Ctx[seatSpec1.Via.Name]
	// and
	recvSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: providerProc.PID,
		Term: step.RecvSpec{
			X: providerProcPID,
			Y: clientProc.PID,
			Cont: step.CloseSpec{
				A: providerProcPID,
			},
		},
	}
	// when
	err = dealApi.Take(recvSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	sendSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: clientProc.PID,
		Term: step.SendSpec{
			A: providerProcPID,
			B: clientProc.PID,
		},
	}
	// and
	err = dealApi.Take(sendSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func TestTakeCaseLab(t *testing.T) {
	// given
	label := core.Label("label-1")
	// and
	roleSpec1 := role.RoleSpec{
		FQN: "role-1",
		St: state.WithSpec{
			Choices: map[core.Label]state.Spec{
				label: state.OneSpec{},
			},
		},
	}
	roleRoot1, err := roleApi.Create(roleSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	roleSpec2 := role.RoleSpec{
		FQN: "role-2",
		St:  state.OneSpec{},
	}
	roleRoot2, err := roleApi.Create(roleSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: roleRoot1.St.RID(),
			St:   roleRoot1.St,
		},
	}
	seatRoot1, err := seatApi.Create(seatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: roleRoot2.St.RID(),
			St:   roleRoot2.St,
		},
		Ctx: []chnl.Spec{
			seatSpec1.Via,
		},
	}
	seatRoot2, err := seatApi.Create(seatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	providerProc, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: providerProc.PID,
		},
	}
	clientProc, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerProcPID := clientProc.Ctx[seatSpec1.Via.Name]
	// and
	caseSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: providerProc.PID,
		Term: step.CaseSpec{
			Z: providerProcPID,
			Conts: map[core.Label]step.Term{
				label: step.CloseSpec{
					A: providerProcPID,
				},
			},
		},
	}
	// when
	err = dealApi.Take(caseSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	labSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: clientProc.PID,
		Term: step.LabSpec{
			C: providerProcPID,
			L: label,
		},
	}
	// and
	err = dealApi.Take(labSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func TestTakeSpawn(t *testing.T) {
	// given
	roleSpec1 := role.RoleSpec{
		FQN: "role-1",
		St:  state.OneSpec{},
	}
	roleRoot1, err := roleApi.Create(roleSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	roleSpec2 := role.RoleSpec{
		FQN: "role-2",
		St:  state.OneSpec{},
	}
	roleRoot2, err := roleApi.Create(roleSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: roleRoot1.St.RID(),
			St:   roleRoot1.St,
		},
	}
	seatRoot1, err := seatApi.Create(seatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: roleRoot2.St.RID(),
			St:   roleRoot2.St,
		},
		Ctx: []chnl.Spec{
			seatSpec1.Via,
		},
	}
	seatRoot2, err := seatApi.Create(seatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	providerProc, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	z := sym.New("foo")
	// and
	spawnSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: providerProc.PID,
		Term: step.SpawnSpec{
			Z: z,
			Ctx: map[chnl.Name]chnl.ID{
				seatSpec1.Via.Name: providerProc.PID,
			},
			Cont: step.CloseSpec{
				A: z,
			},
			SeatID: seatRoot2.ID,
		},
	}
	// when
	err = dealApi.Take(spawnSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func TestTakeFwd(t *testing.T) {
	// given
	oneSpec := role.RoleSpec{
		FQN: "role-1",
		St:  state.OneSpec{},
	}
	oneRoot, err := roleApi.Create(oneSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: oneRoot.St.RID(),
			St:   oneRoot.St,
		},
	}
	seatRoot1, err := seatApi.Create(seatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: oneRoot.St.RID(),
			St:   oneRoot.St,
		},
	}
	seatRoot2, err := seatApi.Create(seatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	seatSpec3 := seat.SeatSpec{
		Name: "seat-3",
		Via: chnl.Spec{
			Name: "chnl-3",
			StID: oneRoot.St.RID(),
			St:   oneRoot.St,
		},
	}
	seatRoot3, err := seatApi.Create(seatSpec3)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	dealRoot, err := dealApi.Create(dealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	providerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	provider1, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	provider2, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	middlemanSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
	}
	middleman, err := dealApi.Involve(middlemanSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot3.ID,
	}
	client, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// when
	fwdSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: middleman.PID,
		Term: step.FwdSpec{
			C: provider1.PID,
			D: provider2.PID,
		},
	}
	// and
	err = dealApi.Take(fwdSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: provider2.PID,
		Term: step.CloseSpec{
			A: provider2.PID,
		},
	}
	// and
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: client.PID,
		Term: step.WaitSpec{
			X: provider1.PID,
			Cont: step.CloseSpec{
				A: middleman.PID,
			},
		},
	}
	// and
	err = dealApi.Take(waitSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}
