package dbmodel

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type proxyType int16

const (
	// ResidentialProxyType резидентские прокси, большинство операций через них
	ResidentialProxyType proxyType = 1
	// CheapProxyType дешёвые прокси, используются для загрузок фотографий постов
	CheapProxyType proxyType = 2
)

// Value is a implementation for driver.Valuer.
func (p Proxy) Value() (driver.Value, error) {
	bytes, err := json.Marshal(p)
	return string(bytes), err
}

// Value is a implementation for driver.Valuer.
func (p Proxy) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

// Scan is an implementation for sql.Scanner.
func (p *Proxy) Scan(value interface{}) error {
	if value == nil {
		*p = Proxy{}
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed for %T = %v", value, value)
	}

	return json.Unmarshal(data, p)
}

func (p Proxy) GetID() uuid.UUID {
	return p.ID
}

func (p Proxy) PythonString() string {
	return fmt.Sprintf("http://%s:%s@%s:%d", p.Login, p.Pass, p.Host, p.Port)
}
