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
		FQN: "seat-1",
		PE: chnl.Spec{
			Name: "chnl-1",
			StID: oneRole.St.Ident(),
		},
	}
	oneSeat1, err := seatApi.Create(oneSeatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec2 := seat.SeatSpec{
		FQN: "seat-2",
		PE: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.Ident(),
		},
		CEs: []chnl.Spec{
			oneSeat1.PE,
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
		Deal: bigDeal.ID,
		Decl: oneSeat1.ID,
	}
	closer, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat2.ID,
		TEs: []chnl.ID{
			closer.ID,
		},
	}
	waiter, err := dealApi.Involve(waiterSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: closer.ID,
		Term: step.CloseSpec{
			A: closer.ID,
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
		PID: waiter.ID,
		Term: step.WaitSpec{
			X: closer.ID,
			Cont: step.CloseSpec{
				A: waiter.ID,
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
		FQN: "seat-1",
		PE: chnl.Spec{
			Name: "chnl-1",
			StID: lolliRole.St.Ident(),
		},
	}
	lolliSeat, err := seatApi.Create(lolliSeatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec1 := seat.SeatSpec{
		FQN: "seat-2",
		PE: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.Ident(),
		},
	}
	oneSeat1, err := seatApi.Create(oneSeatSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec2 := seat.SeatSpec{
		FQN: "seat-3",
		PE: chnl.Spec{
			Name: "chnl-3",
			StID: oneRole.St.Ident(),
		},
		CEs: []chnl.Spec{
			lolliSeatSpec.PE,
			oneSeat1.PE,
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
		Deal: bigDeal.ID,
		Decl: lolliSeat.ID,
	}
	receiver, err := dealApi.Involve(receiverSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	messageSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat1.ID,
	}
	message, err := dealApi.Involve(messageSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	senderSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat2.ID,
		TEs: []chnl.ID{
			receiver.ID,
			message.ID,
		},
	}
	sender, err := dealApi.Involve(senderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	recvSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: receiver.ID,
		Term: step.RecvSpec{
			X: receiver.ID,
			Y: message.ID,
			Cont: step.WaitSpec{
				X: message.ID,
				Cont: step.CloseSpec{
					A: receiver.ID,
				},
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
		PID: sender.ID,
		Term: step.SendSpec{
			A: receiver.ID,
			B: message.ID,
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
		FQN: "seat-1",
		PE: chnl.Spec{
			Name: "chnl-1",
			StID: withRole.St.Ident(),
		},
	}
	withSeat, err := seatApi.Create(withSeatSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeatSpec := seat.SeatSpec{
		FQN: "seat-2",
		PE: chnl.Spec{
			Name: "chnl-2",
			StID: oneRole.St.Ident(),
		},
		CEs: []chnl.Spec{
			withSeat.PE,
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
		Deal: bigDeal.ID,
		Decl: withSeat.ID,
	}
	follower, err := dealApi.Involve(followerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	deciderSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat.ID,
		TEs: []chnl.ID{
			follower.ID,
		},
	}
	decider, err := dealApi.Involve(deciderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	caseSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: follower.ID,
		Term: step.CaseSpec{
			X: follower.ID,
			Conts: map[core.Label]step.Term{
				label: step.CloseSpec{
					A: follower.ID,
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
		PID: decider.ID,
		Term: step.LabSpec{
			A: follower.ID,
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
			FQN: "seat-1",
			PE: chnl.Spec{
				Name: "chnl-1",
				StID: oneRole.St.Ident(),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat2, err := seatApi.Create(
		seat.SeatSpec{
			FQN: "seat-2",
			PE: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				oneSeat1.PE,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat3, err := seatApi.Create(
		seat.SeatSpec{
			FQN: "seat-3",
			PE: chnl.Spec{
				Name: "chnl-3",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				oneSeat1.PE,
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
			Deal: bigDeal.ID,
			Decl: oneSeat1.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	spawnerSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat2.ID,
		TEs: []chnl.ID{
			injectee.ID,
		},
	}
	spawner, err := dealApi.Involve(spawnerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	z := sym.New("z")
	// and
	spawnSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: spawner.ID,
		Term: step.SpawnSpec{
			PE: z,
			CEs: []chnl.ID{
				injectee.ID,
			},
			Cont: step.WaitSpec{
				X: z,
				Cont: step.CloseSpec{
					A: spawner.ID,
				},
			},
			Seat: oneSeat3.ID,
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
			FQN: "seat-1",
			PE: chnl.Spec{
				Name: "chnl-1",
				StID: oneRole.St.Ident(),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat2, err := seatApi.Create(
		seat.SeatSpec{
			FQN: "seat-2",
			PE: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				oneSeat1.PE,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	// and
	oneSeat3, err := seatApi.Create(
		seat.SeatSpec{
			FQN: "seat-3",
			PE: chnl.Spec{
				Name: "chnl-3",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				oneSeat1.PE,
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
		Deal: bigDeal.ID,
		Decl: oneSeat1.ID,
	}
	closer, err := dealApi.Involve(closerSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	forwarderSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat2.ID,
		TEs: []chnl.ID{
			closer.ID,
		},
	}
	forwarder, err := dealApi.Involve(forwarderSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waiterSpec := deal.PartSpec{
		Deal: bigDeal.ID,
		Decl: oneSeat3.ID,
		TEs: []chnl.ID{
			forwarder.ID,
		},
	}
	waiter, err := dealApi.Involve(waiterSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: closer.ID,
		Term: step.CloseSpec{
			A: closer.ID,
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
		PID: forwarder.ID,
		Term: step.FwdSpec{
			C: forwarder.ID,
			D: closer.ID,
		},
	}
	err = dealApi.Take(fwdSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DID: bigDeal.ID,
		PID: waiter.ID,
		Term: step.WaitSpec{
			X: forwarder.ID,
			Cont: step.CloseSpec{
				A: waiter.ID,
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
