package entrance

type ChallengeMethod string

const (
	ChallengeMethodOneTimePasscode ChallengeMethod = "oob"
	ChallengeMethodPassword        ChallengeMethod = "password"
	ChallengeMethodRedirect        ChallengeMethod = "redirect"
)

type Capability string

const (
	CapabilityMultifactorRequired  Capability = "mfa_required"
	CapabilityRegistrationRequired Capability = "registration_required"
)

type AuthenticationMethod struct {
	Identifier       string          `json:"id"`
	ChallengeType    ChallengeMethod `json:"challenge_type"`
	ChallengeChannel string          `json:"challenge_channel"`
	LoginHint        string          `json:"login_hint"`
}

type ChallengeResponse struct {
	ContinuationToken    string          `json:"continuation_token"`
	ChallengeType        ChallengeMethod `json:"challenge_type"`
	ChallengeChannel     string          `json:"challenge_channel"`
	CodeLength           int             `json:"code_length"`
	ChallengeTargetLabel string          `json:"challenge_target_label"`
}

type InitiateResponse struct {
	ContinuationToken string `json:"continuation_token"`
}

type IntrospectResponse struct {
	ContinuationToken string                 `json:"continuation_token"`
	Methods           []AuthenticationMethod `json:"methods"`
}

type RegistrationChallengeResponse struct {
	ContinuationToken    string          `json:"continuation_token"`
	ChallengeType        ChallengeMethod `json:"challenge_type"`
	ChallengeChannel     string          `json:"challenge_channel"`
	CodeLength           int             `json:"code_length"`
	ChallengeTargetLabel string          `json:"challenge_target_label"`
}

type RegistrationContinueResponse struct {
	ContinuationToken string `json:"continuation_token"`
}

type RegistrationIntrospectResponse struct {
	ContinuationToken string                 `json:"continuation_token"`
	Methods           []AuthenticationMethod `json:"methods"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
