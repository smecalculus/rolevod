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
	roleAPI = role.NewAPI()
	sigAPI  = sig.NewAPI()
	dealAPI = deal.NewAPI()
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
	tables := []string{
		"aliases",
		"team_roots",
		"sig_roots", "sig_pes", "sig_ces",
		"role_roots", "role_states",
		"states", "channels", "steps", "clientships"}
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
	ps := deal.Spec{Name: "parent-deal"}
	pr, err := dealAPI.Create(ps)
	if err != nil {
		t.Fatal(err)
	}
	// and
	cs := deal.Spec{Name: "child-deal"}
	cr, err := dealAPI.Create(cs)
	if err != nil {
		t.Fatal(err)
	}
	// when
	ks := deal.KinshipSpec{
		ParentID: pr.ID,
		ChildIDs: []id.ADT{cr.ID},
	}
	err = dealAPI.Establish(ks)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := dealAPI.Retrieve(pr.ID)
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
		oneRoleSpec := role.Spec{
			FQN:   "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerSigSpec := sig.Spec{
			FQN: "closer",
			PE: chnl.Spec{
				Key:  "closing-1",
				Link: oneRole.FQN,
			},
		}
		closerSig, err := sigAPI.Create(closerSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSigSpec := sig.Spec{
			FQN: "waiter",
			PE: chnl.Spec{
				Key:  "closing-2",
				Link: oneRole.FQN,
			},
			CEs: []chnl.Spec{
				closerSig.PE,
			},
		}
		waiterSig, err := sigAPI.Create(waiterSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		bigDealSpec := deal.Spec{
			Name: "big-deal",
		}
		bigDeal, err := dealAPI.Create(bigDealSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  closerSig.ID,
		}
		closer, err := dealAPI.Involve(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  waiterSig.ID,
			TEs: []chnl.ID{
				closer.ID,
			},
		}
		waiter, err := dealAPI.Involve(waiterSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closeSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  closer.ID,
			Term: step.CloseSpec{
				A: closer.ID,
			},
		}
		// when
		err = dealAPI.Take(closeSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  waiter.ID,
			Term: step.WaitSpec{
				X: closer.ID,
				Cont: step.CloseSpec{
					A: waiter.ID,
				},
			},
		}
		// and
		err = dealAPI.Take(waitSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("RecvSend", func(t *testing.T) {
		tc.Setup(t)
		// given
		lolliRoleSpec := role.Spec{
			FQN: "lolli-role",
			State: state.LolliSpec{
				Y: state.OneSpec{},
				Z: state.OneSpec{},
			},
		}
		lolliRole, err := roleAPI.Create(lolliRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := role.Spec{
			FQN:   "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		lolliSigSpec := sig.Spec{
			FQN: "sig-1",
			PE: chnl.Spec{
				Key:  "chnl-1",
				Link: lolliRole.FQN,
			},
		}
		lolliSig, err := sigAPI.Create(lolliSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec1 := sig.Spec{
			FQN: "sig-2",
			PE: chnl.Spec{
				Key:  "chnl-2",
				Link: oneRole.FQN,
			},
		}
		oneSig1, err := sigAPI.Create(oneSigSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec2 := sig.Spec{
			FQN: "sig-3",
			PE: chnl.Spec{
				Key:  "chnl-3",
				Link: oneRole.FQN,
			},
			CEs: []chnl.Spec{
				lolliSigSpec.PE,
				oneSig1.PE,
			},
		}
		oneSig2, err := sigAPI.Create(oneSigSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// and
		bigDealSpec := deal.Spec{
			Name: "deal-1",
		}
		bigDeal, err := dealAPI.Create(bigDealSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		receiverSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  lolliSig.ID,
		}
		receiver, err := dealAPI.Involve(receiverSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		messageSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig1.ID,
		}
		message, err := dealAPI.Involve(messageSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		senderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig2.ID,
			TEs: []chnl.ID{
				receiver.ID,
				message.ID,
			},
		}
		sender, err := dealAPI.Involve(senderSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		recvSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  receiver.ID,
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
		err = dealAPI.Take(recvSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		sendSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  sender.ID,
			Term: step.SendSpec{
				A: receiver.ID,
				B: message.ID,
			},
		}
		// and
		err = dealAPI.Take(sendSpec)
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
		withRoleSpec := role.Spec{
			FQN: "with-role",
			State: state.WithSpec{
				Choices: map[core.Label]state.Spec{
					label: state.OneSpec{},
				},
			},
		}
		withRole, err := roleAPI.Create(withRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneRoleSpec := role.Spec{
			FQN:   "one-role",
			State: state.OneSpec{},
		}
		oneRole, err := roleAPI.Create(oneRoleSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		withSigSpec := sig.Spec{
			FQN: "sig-1",
			PE: chnl.Spec{
				Key:  "chnl-1",
				Link: withRole.FQN,
			},
		}
		withSig, err := sigAPI.Create(withSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSigSpec := sig.Spec{
			FQN: "sig-2",
			PE: chnl.Spec{
				Key:  "chnl-2",
				Link: oneRole.FQN,
			},
			CEs: []chnl.Spec{
				withSig.PE,
			},
		}
		oneSig, err := sigAPI.Create(oneSigSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		bigDealSpec := deal.Spec{
			Name: "deal-1",
		}
		bigDeal, err := dealAPI.Create(bigDealSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		followerSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  withSig.ID,
		}
		follower, err := dealAPI.Involve(followerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		deciderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig.ID,
			TEs: []chnl.ID{
				follower.ID,
			},
		}
		decider, err := dealAPI.Involve(deciderSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		caseSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  follower.ID,
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
		err = dealAPI.Take(caseSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		labSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  decider.ID,
			Term: step.LabSpec{
				A: follower.ID,
				L: label,
			},
		}
		// and
		err = dealAPI.Take(labSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Spawn", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := roleAPI.Create(
			role.Spec{
				FQN:   "one-role",
				State: state.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-1",
				PE: chnl.Spec{
					Key:  "chnl-1",
					Link: oneRole.FQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig2, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-2",
				PE: chnl.Spec{
					Key:  "chnl-2",
					Link: oneRole.FQN,
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
		oneSig3, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-3",
				PE: chnl.Spec{
					Key:  "chnl-3",
					Link: oneRole.FQN,
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
		bigDeal, err := dealAPI.Create(
			deal.Spec{
				Name: "deal-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		injectee, err := dealAPI.Involve(
			deal.PartSpec{
				Deal: bigDeal.ID,
				Sig:  oneSig1.ID,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		spawnerSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig2.ID,
			TEs: []chnl.ID{
				injectee.ID,
			},
		}
		spawner, err := dealAPI.Involve(spawnerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		z := sym.New("z")
		// and
		spawnSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  spawner.ID,
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
		err = dealAPI.Take(spawnSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})

	t.Run("Fwd", func(t *testing.T) {
		tc.Setup(t)
		// given
		oneRole, err := roleAPI.Create(
			role.Spec{
				FQN:   "one-role",
				State: state.OneSpec{},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig1, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-1",
				PE: chnl.Spec{
					Key:  "chnl-1",
					Link: oneRole.FQN,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		oneSig2, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-2",
				PE: chnl.Spec{
					Key:  "chnl-2",
					Link: oneRole.FQN,
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
		oneSig3, err := sigAPI.Create(
			sig.Spec{
				FQN: "sig-3",
				PE: chnl.Spec{
					Key:  "chnl-3",
					Link: oneRole.FQN,
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
		bigDeal, err := dealAPI.Create(
			deal.Spec{
				Name: "deal-1",
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closerSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig1.ID,
		}
		closer, err := dealAPI.Involve(closerSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		forwarderSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig2.ID,
			TEs: []chnl.ID{
				closer.ID,
			},
		}
		forwarder, err := dealAPI.Involve(forwarderSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waiterSpec := deal.PartSpec{
			Deal: bigDeal.ID,
			Sig:  oneSig3.ID,
			TEs: []chnl.ID{
				forwarder.ID,
			},
		}
		waiter, err := dealAPI.Involve(waiterSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		closeSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  closer.ID,
			Term: step.CloseSpec{
				A: closer.ID,
			},
		}
		err = dealAPI.Take(closeSpec)
		if err != nil {
			t.Fatal(err)
		}
		// when
		fwdSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			// канал пересыльщика должен закрыться?
			PID: forwarder.ID,
			Term: step.FwdSpec{
				C: forwarder.ID,
				D: closer.ID,
			},
		}
		err = dealAPI.Take(fwdSpec)
		if err != nil {
			t.Fatal(err)
		}
		// and
		waitSpec := deal.TranSpec{
			Deal: bigDeal.ID,
			PID:  waiter.ID,
			Term: step.WaitSpec{
				X: forwarder.ID,
				Cont: step.CloseSpec{
					A: waiter.ID,
				},
			},
		}
		err = dealAPI.Take(waitSpec)
		if err != nil {
			t.Fatal(err)
		}
		// then
		// TODO добавить проверку
	})
}
