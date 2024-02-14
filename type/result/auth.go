package result

type ServiceLoginResult struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
