// Package fixture holds the fakes to test code that receives the answer of
// the NATS connection question, without a NATS server:
//
//	closed := fixture.CreateFakeClosedConnection()
//	fixture.AssertMatch(t, been.Closed(), closed)
package fixture
