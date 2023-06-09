/*		Basic One Time Password - OTP solution for Authorisation
*		- happens before Websocket connection is established.
*		- OTPs expire after a few seconds
*/
package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// holds all active/valid tokens
type RetentionMap map[string]OTP

func NewRetentionMap(ctx context.Context, duration time.Duration) RetentionMap {
	rm := make(RetentionMap)
	go rm.retention(ctx, duration)
	return rm
}

// a single (valid) login token
type OTP struct {
	Key      string // autogenerated key we pass down
	Created  time.Time
	Username string
	Password string
}

func (rm RetentionMap) NewOTP(username, password string) OTP {
	otp := OTP{
		Key:      uuid.NewString(),
		Created:  time.Now(),
		Username: username,
		Password: password,
	}
	rm[otp.Key] = otp
	return otp
}

// returns (if exists) OTK in retention map and a boolean if it exists
func (rm RetentionMap) VertifyOTP(otp_key string) (OTP, bool) {
	otk, ok := rm[otp_key]
	if !ok {
		return OTP{}, false
	}
	delete(rm, otp_key)
	return otk, true
}

// runs in background and removes expired OTPs (is blocking so async save)
func (rm RetentionMap) retention(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(500 * time.Millisecond) //frequency we check
	for {
		select {
		case <-ticker.C:
			for _, otp := range rm {
				if otp.Created.Add(duration).Before(time.Now()) {
					delete(rm, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
