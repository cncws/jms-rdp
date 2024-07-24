package client

import (
	"net/http"
	"time"

	"gopkg.in/twindagger/httpsig.v1"
)

type sigAuth struct {
	KeyID    string
	SecretID string
}

func (auth *sigAuth) Sign(r *http.Request) error {
	headers := []string{"(request-target)", "date"}
	signer, err := httpsig.NewRequestSigner(auth.KeyID, auth.SecretID, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

type sigAuthRoundTripper struct {
	http.RoundTripper
	SigAuth *sigAuth
}

func (rt *sigAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err := rt.SigAuth.Sign(req); err != nil {
		return nil, err
	}
	return rt.RoundTripper.RoundTrip(req)
}

func NewSigAuthRoundTripper(keyID, secretID string) *sigAuthRoundTripper {
	return &sigAuthRoundTripper{
		RoundTripper: http.DefaultTransport,
		SigAuth: &sigAuth{
			KeyID:    keyID,
			SecretID: secretID,
		},
	}
}
