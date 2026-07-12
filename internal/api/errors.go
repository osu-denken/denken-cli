package api

import "fmt"

// APIError は API が 4xx/5xx を返したときのエラー。
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("api error: status %d", e.Status)
	}
	return fmt.Sprintf("api error: status %d: %s", e.Status, e.Body)
}

// RateLimited はレート制限 (429) かどうかを返す。
func (e *APIError) RateLimited() bool {
	return e.Status == 429
}
