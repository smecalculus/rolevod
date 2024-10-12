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
	oneRoleSpec := role.RoleSpec{
		FQN: "one-role",
		St:  state.OneSpec{},
	}
	oneRole, err := roleApi.Create(oneRoleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec1 := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: oneRole.St.RID(),
			St:   oneRole.St,
		},
	}
	oneSeat1, err := seatApi.Create(oneSeatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec2 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.RID(),
			St:   oneRole.St,
		},
		Ctx: []chnl.Spec{
			oneSeat1.Via,
		},
	}
	oneSeat2, err := seatApi.Create(oneSeatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	bigDealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	bigDeal, err := dealApi.Create(bigDealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat1.ID,
	}
	closer, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			oneSeat1.Via.Name: closer.PID,
		},
	}
	waiter, err := dealApi.Involve(waiterSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: closer.PID,
		Term: step.CloseSpec{
			A: closer.PID,
		},
	}
	// when
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: waiter.PID,
		Term: step.WaitSpec{
			X: closer.PID,
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
	// then
	// TODO добавить проверку
}

func TestTakeRecvSend(t *testing.T) {
	// given
	lolliRoleSpec := role.RoleSpec{
		FQN: "lolli-role",
		St: state.LolliSpec{
			Y: state.OneSpec{},
			Z: state.OneSpec{},
		},
	}
	lolliRole, err := roleApi.Create(lolliRoleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneRoleSpec := role.RoleSpec{
		FQN: "one-role",
		St:  state.OneSpec{},
	}
	oneRole, err := roleApi.Create(oneRoleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	lolliSeatSpec := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: lolliRole.St.RID(),
			St:   lolliRole.St,
		},
	}
	lolliSeat, err := seatApi.Create(lolliSeatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec1 := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.RID(),
			St:   oneRole.St,
		},
	}
	oneSeat1, err := seatApi.Create(oneSeatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec2 := seat.SeatSpec{
		Name: "seat-3",
		Via: chnl.Spec{
			Name: "chnl-3",
			StID: oneRole.St.RID(),
			St:   oneRole.St,
		},
		Ctx: []chnl.Spec{
			lolliSeatSpec.Via,
			oneSeat1.Via,
		},
	}
	oneSeat2, err := seatApi.Create(oneSeatSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	bigDealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	bigDeal, err := dealApi.Create(bigDealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	receiverSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: lolliSeat.ID,
	}
	receiver, err := dealApi.Involve(receiverSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	messageSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat1.ID,
	}
	message, err := dealApi.Involve(messageSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			lolliSeat.Via.Name: receiver.PID,
			oneSeat1.Via.Name:  message.PID,
		},
	}
	sender, err := dealApi.Involve(senderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	recvSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: receiver.PID,
		Term: step.RecvSpec{
			X: receiver.PID,
			Y: message.PID,
			// закрываемся с каналом в контексте
			Cont: step.CloseSpec{
				A: receiver.PID,
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
		DID: bigDeal.ID,
		PID: sender.PID,
		Term: step.SendSpec{
			A: receiver.PID,
			B: message.PID,
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
	withRoleSpec := role.RoleSpec{
		FQN: "with-role",
		St: state.WithSpec{
			Choices: map[core.Label]state.Spec{
				label: state.OneSpec{},
			},
		},
	}
	withRole, err := roleApi.Create(withRoleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneRoleSpec := role.RoleSpec{
		FQN: "one-role",
		St:  state.OneSpec{},
	}
	oneRole, err := roleApi.Create(oneRoleSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	withSeatSpec := seat.SeatSpec{
		Name: "seat-1",
		Via: chnl.Spec{
			Name: "chnl-1",
			StID: withRole.St.RID(),
			St:   withRole.St,
		},
	}
	withSeat, err := seatApi.Create(withSeatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec := seat.SeatSpec{
		Name: "seat-2",
		Via: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.RID(),
			St:   oneRole.St,
		},
		Ctx: []chnl.Spec{
			withSeat.Via,
		},
	}
	oneSeat, err := seatApi.Create(oneSeatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	bigDealSpec := deal.DealSpec{
		Name: "deal-1",
	}
	bigDeal, err := dealApi.Create(bigDealSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	followerSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: withSeat.ID,
	}
	follower, err := dealApi.Involve(followerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	deciderSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat.ID,
		Ctx: map[chnl.Name]chnl.ID{
			withSeat.Via.Name: follower.PID,
		},
	}
	decider, err := dealApi.Involve(deciderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	caseSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: follower.PID,
		Term: step.CaseSpec{
			Z: follower.PID,
			Conts: map[core.Label]step.Term{
				label: step.CloseSpec{
					A: follower.PID,
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
		DID: bigDeal.ID,
		PID: decider.PID,
		Term: step.LabSpec{
			C: follower.PID,
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
			FQN: "one-role",
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
	bigDeal, err := dealApi.Create(
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
			DealID: bigDeal.ID,
			SeatID: oneSeat1.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	spawnerSpec := deal.PartSpec{
		DealID: bigDeal.ID,
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
	spawnSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: spawner.PID,
		Term: step.SpawnSpec{
			Z: sym.New("z"),
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
			FQN: "one-role",
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
	bigDeal, err := dealApi.Create(
		deal.DealSpec{
			Name: "deal-1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closerSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat1.ID,
	}
	closer, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat2.ID,
		Ctx: map[chnl.Name]chnl.ID{
			oneSeat1.Via.Name: closer.PID,
		},
	}
	forwarder, err := dealApi.Involve(forwarderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterSpec := deal.PartSpec{
		DealID: bigDeal.ID,
		SeatID: oneSeat3.ID,
		Ctx: map[chnl.Name]chnl.ID{
			oneSeat1.Via.Name: forwarder.PID,
		},
	}
	waiter, err := dealApi.Involve(waiterSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: closer.PID,
		Term: step.CloseSpec{
			A: closer.PID,
		},
	}
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// when
	fwdSpec := deal.TranSpec{
		DID: bigDeal.ID,
		// канал пересыльщика должен закрыться?
		PID: forwarder.PID,
		Term: step.FwdSpec{
			C: closer.PID,
			D: forwarder.PID,
		},
	}
	err = dealApi.Take(fwdSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: waiter.PID,
		Term: step.WaitSpec{
			X: forwarder.PID,
			Cont: step.CloseSpec{
				A: waiter.PID,
			},
		},
	}
	err = dealApi.Take(waitSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}
