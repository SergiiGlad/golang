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
  log.Debugf("Start generating token for email %s", email)
  user, err := ts.UD.FindUserByEmail(email)
  if err != nil {
    return "", err
  }
  user.AccStatus = entity.InActive
  log.Debugf("User fore email: %s", email)
  token := randString(tokenLength)
  log.Debugf("New token %s", token)
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
  user, err := ts.UD.FindUserById(t.UserId)
  if err != nil || &user == nil {
    return false, err
  }
  user.AccStatus = entity.Active
  _, err = ts.UD.UpdateUser(user.ID, &user)
  if err != nil {
    return false, err
  }
  t.IsActive = false
  err = ts.TD.UpdateToken(t.ID, t)
  if err != nil {
    return false, err
  }
  log.Infof("Successfully approved user %s for token %s", user, t)
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
