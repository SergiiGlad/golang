package controllers

import (
  "math/rand"
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao/entity"
)

type TokenService struct {
  UD interfaces.UserDao
  TD interfaces.UserTokenDao
}

var _ TokenGeneratorInterface = &TokenService{}

const tokenLength = 64

func (ts TokenService) GenerateTokenForEmail(email string) (string, error) {
  log.Debug("Start generating token for email {}", email)
  user, err := ts.UD.FindUserByEmail(email)
  if err != nil {
    return "", err
  }
  user.AccStatus = entity.InActive
  log.Debug("User fore email: {}")
  _, err = ts.UD.UpdateUser(user.ID, &user)
  if err != nil {
    return "", err
  }
  token := randString(tokenLength)
  log.Debug("New token {}", token)
  _, err = ts.TD.AddToken(entity.UserToken{
    Token:    token,
    Email:    email,
    IsActive: true,
    UserId:   user.ID,
  })

  if err != nil {
    return "", err
  }
  return token, nil
}

func (ts TokenService) ApproveUser(token string) (bool, error) {
  t, err := ts.TD.FindTokenByToken(token)
  if err != nil || t == nil || t.Token != token {
    return false, err
  }
  user, err := ts.UD.FindUserByEmail(t.Email)
  if err != nil || &user == nil {
    return false, err
  }
  user.AccStatus = entity.Active
  _, err = ts.UD.UpdateUser(user.ID, &user)
  if err != nil {
    return false, err
  }
  return true, nil
}

const (
  letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  letterIdxBits = 6
  letterIdxMask = 1<<letterIdxBits - 1
  letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
  b := make([]byte, n)
  for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
    if remain == 0 {
      cache, remain = rand.Int63(), letterIdxMax
    }
    if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
      b[i] = letterBytes[idx]
      i--
    }
    cache >>= letterIdxBits
    remain--
  }
  return string(b)
}
