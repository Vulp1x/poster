package headers

type AuthorizationData struct {
	DsUserID  string `json:"ds_user_id"`
	SessionID string `json:"sessionid"` // !!! не трогай иначе сломается парсинг токена!!!
	CSRFToken string `json:"csrf_token"`
}

type Base struct {
	Mid             string
	DsUserID        string
	Rur             string
	Authorization   string
	WWWClaim        string
	AuthData        AuthorizationData
	BlocksVersionID string
}

// X-MID=;
// IG-U-DS-USER-ID=55063899557;
// IG-U-RUR=ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62;
// Authorization=Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9;
// X-IG-WWW-Claim=hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY"

// X-MID=;IG-U-DS-USER-ID=55063899557;IG-U-RUR=ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62;Authorization=Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9;X-IG-WWW-Claim=hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY"
