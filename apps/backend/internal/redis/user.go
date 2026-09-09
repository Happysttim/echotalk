package redis

import (
	"context"
	"strings"
	"time"
)

type VerifyType string

const (
	RegisterVerify VerifyType = "register"
	PasswordVerify VerifyType = "password"
	UnknownVerify  VerifyType = "unknown"
)

func IsVerifyType(s string) bool {
	return strings.Compare(s, string(RegisterVerify)) == 0 ||
		strings.Compare(s, string(PasswordVerify)) == 0
}

func IsRegisterVerify(s string) bool {
	return strings.Compare(s, string(RegisterVerify)) == 0
}

func IsPasswordVerify(s string) bool {
	return strings.Compare(s, string(PasswordVerify)) == 0
}

type RedisUser struct{}

func (r RedisUser) SetVerify(ctx context.Context, verifyLink string, email string, verifyType VerifyType) error {

	if err := RedisClient.Set(ctx, strings.Join([]string{string(verifyType), verifyLink}, ":"), email, time.Minute*3).Err(); err != nil {
		return err
	}

	if err := RedisClient.Set(ctx, strings.Join([]string{string(verifyType), email}, ":"), "verify", time.Second*30).Err(); err != nil {
		return err
	}

	return nil
}

func (r RedisUser) IsVerifyEmail(ctx context.Context, email string, verifyType VerifyType) error {
	_, err := RedisClient.Get(ctx, strings.Join([]string{string(verifyType), email}, ":")).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisUser) IsVerifyLink(ctx context.Context, verifyLink string, verifyType VerifyType) error {
	_, err := RedisClient.Get(ctx, strings.Join([]string{string(verifyType), verifyLink}, ":")).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisUser) DeleteVerify(ctx context.Context, verifyLink string, verifyType VerifyType) error {
	err := RedisClient.Del(ctx, strings.Join([]string{string(verifyType), verifyLink}, ":")).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisUser) GetDelVerify(ctx context.Context, verifyLink string, verifyType VerifyType) (string, error) {
	email, err := RedisClient.GetDel(ctx, strings.Join([]string{string(verifyType), verifyLink}, ":")).Result()
	if err != nil {
		return "", err
	}

	_, err = RedisClient.GetDel(ctx, email).Result()
	if err != nil {
		return "", err
	}

	return email, nil
}
