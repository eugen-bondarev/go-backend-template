package svc

import (
	"testing"
	"time"
)

type testCase struct {
	signing ISigning
	claims  map[string]any
}

func newTestCase() testCase {
	return testCase{
		signing: NewJWTSigning("foobar baz"),
		claims: map[string]any{
			"lorem": "ipsum",
		},
	}
}

func Test_JWTSigning(t *testing.T) {
	tc := newTestCase()

	token, err := tc.signing.Sign(tc.claims, time.Now().Add(time.Millisecond*200))

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(token.Value) == 0 {
		t.Errorf("token is empty")
	}

	parsed, err := tc.signing.Parse(token.Value)

	if parsed["lorem"] != tc.claims["lorem"] {
		t.Errorf("wrong claims. Expected %v, got %v", tc.claims, parsed)
	}
}

func Test_JWTExpiration(t *testing.T) {
	tc := newTestCase()

	token, err := tc.signing.Sign(tc.claims, time.Now().Add(-time.Second*10))

	if err != nil {
		t.Errorf(err.Error())
	}

	parsed, err := tc.signing.Parse(token.Value)

	if err == nil {
		t.Errorf("expected err")
	}

	if parsed["lorem"] == tc.claims["lorem"] {
		t.Errorf("expected err")
	}
}
