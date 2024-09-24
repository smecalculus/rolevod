package deal_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

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
		ParentID:    pr.ID,
		ChildrenIDs: []id.ADT[deal.ID]{cr.ID},
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
		Name: "role-1",
		St:   state.OneSpec{},
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
		Via: chnl.Spec{
			Name: "chnl-2",
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
	partSpec1 := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	chnlRef1, err := dealApi.Involve(partSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	partSpec2 := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
	}
	chnlRef2, err := dealApi.Involve(partSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	waitSpec := deal.TranSpec{
		DealID: dealRoot.ID,
		Term: step.WaitSpec{
			X:    chnlRef1,
			Cont: step.CloseSpec{A: chnlRef2},
		},
	}
	// when
	err = dealApi.Take(waitSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	closeSpec := deal.TranSpec{
		DealID: dealRoot.ID,
		Term: step.CloseSpec{
			A: chnlRef1,
		},
	}
	// and
	err = dealApi.Take(closeSpec)
	if err != nil {
		t.Fatal(err)
	}
	// then
	// TODO добавить проверку
}

func TestTakeRecvSend(t *testing.T) {
	// given
	roleSpec1 := role.RoleSpec{
		Name: "role-1",
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
		Name: "role-2",
		St:   state.OneSpec{},
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
			St:   roleRoot2.St,
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
	partSpec1 := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot1.ID,
	}
	chnlRef1, err := dealApi.Involve(partSpec1)
	if err != nil {
		t.Fatal(err)
	}
	// and
	partSpec2 := deal.PartSpec{
		DealID: dealRoot.ID,
		SeatID: seatRoot2.ID,
	}
	chnlRef2, err := dealApi.Involve(partSpec2)
	if err != nil {
		t.Fatal(err)
	}
	// and
	recvSpec := deal.TranSpec{
		DealID: dealRoot.ID,
		AK:     chnlRef1.PAK,
		Term: step.RecvSpec{
			X:    chnlRef1,
			Y:    chnlRef2,
			Cont: step.CloseSpec{A: chnlRef1},
		},
	}
	// when
	err = dealApi.Take(recvSpec)
	if err != nil {
		t.Fatal(err)
	}
	// and
	sendSpec := deal.TranSpec{
		DealID: dealRoot.ID,
		AK:     chnlRef1.CAK,
		Term: step.SendSpec{
			A: chnlRef1,
			B: chnlRef2,
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
