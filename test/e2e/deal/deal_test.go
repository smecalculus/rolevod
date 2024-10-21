package deal_test

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
	"smecalculus/rolevod/internal/step"

	"smecalculus/rolevod/app/deal"
	"smecalculus/rolevod/app/role"
	"smecalculus/rolevod/app/sig"
)

var (
	roleApi = role.NewRoleApi()
	sigApi  = sig.NewSigApi()
	dealApi = deal.NewDealApi()
	tc      *testCase
)

func TestMain(m *testing.M) {
	ts := testSuite{}
	tc = ts.Setup()
	code := m.Run()
	ts.Teardown()
	os.Exit(code)
}

type testSuite struct {
	db *sql.DB
}

func (ts *testSuite) Setup() *testCase {
	db, err := sql.Open("pgx", "postgres://rolevod:rolevod@localhost:5432/rolevod")
	if err != nil {
		panic(err)
	}
	ts.db = db
	return &testCase{db}
}

func (ts *testSuite) Teardown() {
	err := ts.db.Close()
	if err != nil {
		panic(err)
	}
}

type testCase struct {
	db *sql.DB
}

func (tc *testCase) Setup(t *testing.T) {
	tables := []string{"aliases", "roles", "signatures", "states", "channels", "steps", "clientships"}
	for _, table := range tables {
		_, err := tc.db.Exec(fmt.Sprintf("truncate table %v", table))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEstablishKinship(t *testing.T) {
	tc.Setup(t)
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

func TestTake(t *testing.T) {

	t.Run("WaitClose", func(t *testing.T) {
		tc.Setup(t)
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
		oneSigSpec1 := sig.Spec{
			FQN: "sig-1",
			PE: chnl.Spec{
				Name: "chnl-1",
				StID: oneRole.St.Ident(),
			},
		}
		oneSig1, err := sigApi.Create(oneSigSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec2 := sig.Spec{
			FQN: "sig-2",
			PE: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				oneSig1.PE,
			},
		}
		oneSig2, err := sigApi.Create(oneSigSpec2)
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
			Decl: oneSig1.ID,
		}
		closer, err := dealApi.Involve(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig2.ID,
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
	})

	t.Run("RecvSend", func(t *testing.T) {
		tc.Setup(t)
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
		lolliSigSpec := sig.Spec{
			FQN: "sig-1",
			PE: chnl.Spec{
				Name: "chnl-1",
				StID: lolliRole.St.Ident(),
			},
		}
		lolliSig, err := sigApi.Create(lolliSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec1 := sig.Spec{
			FQN: "sig-2",
			PE: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.Ident(),
			},
		}
		oneSig1, err := sigApi.Create(oneSigSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec2 := sig.Spec{
			FQN: "sig-3",
			PE: chnl.Spec{
				Name: "chnl-3",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				lolliSigSpec.PE,
				oneSig1.PE,
			},
		}
		oneSig2, err := sigApi.Create(oneSigSpec2)
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
			Decl: lolliSig.ID,
		}
		receiver, err := dealApi.Involve(receiverSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		messageSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig1.ID,
		}
		message, err := dealApi.Involve(messageSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		senderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig2.ID,
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
	})

	t.Run("CaseLab", func(t *testing.T) {
		tc.Setup(t)
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
		withSigSpec := sig.Spec{
			FQN: "sig-1",
			PE: chnl.Spec{
				Name: "chnl-1",
				StID: withRole.St.Ident(),
			},
		}
		withSig, err := sigApi.Create(withSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec := sig.Spec{
			FQN: "sig-2",
			PE: chnl.Spec{
				Name: "chnl-2",
				StID: oneRole.St.Ident(),
			},
			CEs: []chnl.Spec{
				withSig.PE,
			},
		}
		oneSig, err := sigApi.Create(oneSigSpec)
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
			Decl: withSig.ID,
		}
		follower, err := dealApi.Involve(followerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		deciderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig.ID,
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
	})

	t.Run("Spawn", func(t *testing.T) {
		tc.Setup(t)
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
		oneSig1, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-1",
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
		oneSig2, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-2",
				PE: chnl.Spec{
					Name: "chnl-2",
					StID: oneRole.St.Ident(),
				},
				CEs: []chnl.Spec{
					oneSig1.PE,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig3, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-3",
				PE: chnl.Spec{
					Name: "chnl-3",
					StID: oneRole.St.Ident(),
				},
				CEs: []chnl.Spec{
					oneSig1.PE,
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
				Decl: oneSig1.ID,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		spawnerSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig2.ID,
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
				Sig: oneSig3.ID,
			},
		}
		// when
		err = dealApi.Take(spawnSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Fwd", func(t *testing.T) {
		tc.Setup(t)
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
		oneSig1, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-1",
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
		oneSig2, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-2",
				PE: chnl.Spec{
					Name: "chnl-2",
					StID: oneRole.St.Ident(),
				},
				CEs: []chnl.Spec{
					oneSig1.PE,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig3, err := sigApi.Create(
			sig.Spec{
				FQN: "sig-3",
				PE: chnl.Spec{
					Name: "chnl-3",
					StID: oneRole.St.Ident(),
				},
				CEs: []chnl.Spec{
					oneSig1.PE,
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
			Decl: oneSig1.ID,
		}
		closer, err := dealApi.Involve(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		forwarderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Decl: oneSig2.ID,
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
			Decl: oneSig3.ID,
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
	})
}
