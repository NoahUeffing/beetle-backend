package healthcheck

type Status struct {
	Name     string   `json:"name"`
	Up       bool     `json:"up"`
	Messages []string `json:"messages"`
}

type IHealthCheckService interface {
	Check() Status
}
