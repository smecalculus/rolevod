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
		Ctx: []chnl.Spec{
			seatSpec1.Via,
		},
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
	provider, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: provider.PID,
		},
	}
	client, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: provider.PID,
		Term: step.CloseSpec{
			A: provider.PID,
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
		PID: client.PID,
		Term: step.WaitSpec{
			X: provider.PID,
			Cont: step.CloseSpec{
				A: client.PID,
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
	lolliSpec := role.RoleSpec{
		FQN: "role-1",
		St: state.LolliSpec{
			Y: state.OneSpec{},
			Z: state.OneSpec{},
		},
	}
	lolliRoot, err := roleApi.Create(lolliSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSpec := role.RoleSpec{
		FQN: "role-2",
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
			StID: lolliRoot.St.RID(),
			St:   lolliRoot.St,
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
	provider, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: provider.PID,
		},
	}
	client, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	recvSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: provider.PID,
		Term: step.RecvSpec{
			X: provider.PID,
			Y: client.PID,
			Cont: step.CloseSpec{
				A: provider.PID,
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
		PID: client.PID,
		Term: step.SendSpec{
			A: provider.PID,
			// отправляет сам себя!
			B: client.PID,
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
	withSpec := role.RoleSpec{
		FQN: "role-1",
		St: state.WithSpec{
			Choices: map[core.Label]state.Spec{
				label: state.OneSpec{},
			},
		},
	}
	withRoot, err := roleApi.Create(withSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSpec := role.RoleSpec{
		FQN: "role-2",
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
			StID: withRoot.St.RID(),
			St:   withRoot.St,
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
	provider, err := dealApi.Involve(providerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	clientSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			seatSpec1.Via.Name: provider.PID,
		},
	}
	client, err := dealApi.Involve(clientSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	caseSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: provider.PID,
		Term: step.CaseSpec{
			Z: provider.PID,
			Conts: map[core.Label]step.Term{
				label: step.CloseSpec{
					A: provider.PID,
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
		PID: client.PID,
		Term: step.LabSpec{
			C: provider.PID,
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
	oneRole, err := roleApi.Create(
		role.RoleSpec{
			FQN: "role-1",
			St:  state.OneSpec{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat1, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-1",
			Via: chnl.Spec{
				Name: "chnl-1",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat2, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-2",
			Via: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
			Ctx: []chnl.Spec{
				oneSeat1.Via,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat3, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-3",
			Via: chnl.Spec{
				Name: "chnl-3",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
			Ctx: []chnl.Spec{
				oneSeat1.Via,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealRoot, err := dealApi.Create(
		deal.DealSpec{
			Name: "deal-1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	injectee, err := dealApi.Involve(
		deal.PartSpec{
			DealID: dealRoot.ID,
			SeatID: oneSeat1.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	spawnerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: oneSeat2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			oneSeat1.Via.Name: injectee.PID,
		},
	}
	spawner, err := dealApi.Involve(spawnerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	z := sym.New("foo")
	// and
	spawnSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: spawner.PID,
		Term: step.SpawnSpec{
			Z: z,
			Ctx: map[chnl.Name]chnl.ID{
				oneSeat1.Via.Name: injectee.PID,
			},
			Cont: step.CloseSpec{
				A: spawner.PID,
			},
			SeatID: oneSeat3.ID,
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
	oneRole, err := roleApi.Create(
		role.RoleSpec{
			FQN: "role-1",
			St:  state.OneSpec{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat1, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-1",
			Via: chnl.Spec{
				Name: "chnl-1",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat2, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-2",
			Via: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
			Ctx: []chnl.Spec{
				oneSeat1.Via,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat3, err := seatApi.Create(
		seat.SeatSpec{
			Name: "seat-3",
			Via: chnl.Spec{
				Name: "chnl-3",
				StID: oneRole.St.RID(),
				St:   oneRole.St,
			},
			Ctx: []chnl.Spec{
				oneSeat1.Via,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	dealRoot, err := dealApi.Create(
		deal.DealSpec{
			Name: "deal-1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: oneSeat1.ID,
	}
	closer1, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	closer2, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiter, err := dealApi.Involve(
		deal.PartSpec{
			DealID: dealRoot.ID,
			SeatID: oneSeat2.ID,
			Ctx: map[chnl.Name]chnl.ID{
				oneSeat1.Via.Name: closer1.PID,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderSpec := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: oneSeat3.ID,
		Ctx: map[chnl.Name]chnl.ID{
			oneSeat1.Via.Name: closer1.PID,
		},
	}
	forwarder, err := dealApi.Involve(forwarderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: waiter.PID,
		Term: step.WaitSpec{
			X: closer1.PID,
			Cont: step.CloseSpec{
				A: waiter.PID,
			},
		},
	}
	// and
	err = dealApi.Take(waitSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: dealRoot.ID,
		PID: closer2.PID,
		Term: step.CloseSpec{
			A: closer2.PID,
		},
	}
	// and
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// when
	fwdSpec := deal.TranSpec{
		DID: dealRoot.ID,
		// канал пересыльщика должен закрыться?
		PID: forwarder.PID,
		Term: step.FwdSpec{
			C: closer1.PID,
			D: closer2.PID,
		},
	}
	// and
	err = dealApi.Take(fwdSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}
