package controllers

import (
  "math/rand"
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao"
)

type TokenService struct {
  userDao  interfaces.UserDao
  tokenDao interfaces.UserTokenDao
}

var _ TokenGeneratorInterface = &TokenService{}

const tokenLength = 64

func (ts TokenService) GenerateTokenForEmail(email string) (string, error) {
  user, err := ts.userDao.FindUserByEmail(email)
  if err != nil {
    return "", err
  }
  user.AccStatus = dao.InActive
  err = ts.userDao.UpdateUser(user.ID, user)
  if err != nil {
    return "", err
  }
  token := randString(tokenLength)
  _, err = ts.tokenDao.AddToken(dao.UserToken{
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

func (ts TokenService) ApproveUser(email string, token string) (bool, error) {
  t, err := ts.tokenDao.FindTokenByToken(token)
  if err != nil || t == nil || t.Token != token {
    return false, err
  }
  user, err := ts.userDao.FindUserByEmail(email)
  if err != nil || user == nil {
    return false, err
  }
  user.AccStatus = dao.Active
  err = ts.userDao.UpdateUser(user.ID, user)
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
