package sipgo

import (
	"testing"

	"github.com/emiago/sipgo/sip"
)

// A CANCEL goes where its INVITE went (RFC 3261 §9.1). An INVITE dialled at
// an explicit destination -- an outbound proxy named without a Route header
// -- must produce a CANCEL with that same destination; before this test the
// CANCEL fell back to the Request-URI and was sent to the far end directly.
func TestNewCancelRequestKeepsInviteDestination(t *testing.T) {
	inv := newTestInvite("z9hG4bK.test", "cancel-dest@test")
	inv.SetDestination("10.1.3.24:5062")

	cancel := newCancelRequest(inv)
	if got := cancel.Destination(); got != "10.1.3.24:5062" {
		t.Fatalf("CANCEL destination = %q, want the INVITE's %q", got, "10.1.3.24:5062")
	}
	if cancel.Recipient.String() != inv.Recipient.String() {
		t.Fatalf("CANCEL Request-URI = %s, want the INVITE's %s", cancel.Recipient.String(), inv.Recipient.String())
	}

	// Without an explicit destination the CANCEL resolves like the INVITE
	// did: Route first, else the Request-URI.
	plain := newTestInvite("z9hG4bK.plain", "cancel-plain@test")
	if got := newCancelRequest(plain).Destination(); got != "203.0.113.7:5060" {
		t.Fatalf("plain CANCEL destination = %q, want the Request-URI host:port", got)
	}
}

func newTestInvite(branch, callID string) *sip.Request {
	inv := sip.NewRequest(sip.INVITE, sip.Uri{User: "alice", Host: "203.0.113.7", Port: 5060})
	via := sip.ViaHeader{ProtocolName: "SIP", ProtocolVersion: "2.0", Transport: "UDP", Host: "10.1.1.21", Port: 5080, Params: sip.NewParams()}
	via.Params.Add("branch", branch)
	from := sip.FromHeader{Address: sip.Uri{User: "bob", Host: "10.1.1.21"}, Params: sip.NewParams()}
	from.Params.Add("tag", "ft")
	inv.AppendHeader(&via)
	inv.AppendHeader(&from)
	inv.AppendHeader(&sip.ToHeader{Address: sip.Uri{User: "alice", Host: "203.0.113.7"}})
	inv.AppendHeader(sip.NewHeader("Call-ID", callID))
	return inv
}
