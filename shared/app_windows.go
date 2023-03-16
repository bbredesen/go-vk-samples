//go:build windows

package shared

func NewApp() (App, error) {
	return nil, errors.New("windows not yet implemented")
}
