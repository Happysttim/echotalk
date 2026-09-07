package redis

import (
	"context"
	"time"
)

type RedisUser struct{}

func (r RedisUser) SetVerify(ctx context.Context, verifyLink string, email string) error {
	if err := RedisClient.Set(ctx, verifyLink, email, time.Minute*3).Err(); err != nil {
		return err
	}

	if err := RedisClient.Set(ctx, email, "verify", time.Minute*3).Err(); err != nil {
		return err
	}

	return nil
}

func (r RedisUser) IsVerify(ctx context.Context, email string) error {
	_, err := RedisClient.Get(ctx, email).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisUser) DeleteVerify(ctx context.Context, verifyLink string) error {
	err := RedisClient.Del(ctx, verifyLink).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisUser) GetDelVerify(ctx context.Context, verifyLink string) (string, error) {
	email, err := RedisClient.GetDel(ctx, verifyLink).Result()
	if err != nil {
		return "", err
	}

	_, err = RedisClient.GetDel(ctx, email).Result()
	if err != nil {
		return "", err
	}

	return email, nil
}
