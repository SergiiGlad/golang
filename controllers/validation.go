package controllers

import "regexp"

func ValidRegexItem(item string, pattern string) bool {
  itemRegex := regexp.MustCompile(pattern)

  if isItemOk := itemRegex.MatchString(item); isItemOk == false {
    return false
  }

  return true
}

func ValidEmail(email string) bool {
  return ValidRegexItem(email, "^[a-z0-9]+@[a-z]+[.][a-z]+$")
}

func ValidPhone(phone string) bool {
  return ValidRegexItem(phone, "^[+][0-9]{12}")
}

func ValidCyrillicName(name string) bool {
  return ValidRegexItem(name, "^[А-Я][а-я]{1,49}")
}

func ValidLatinName(name string) bool {
  return ValidRegexItem(name, "^[A-Z][a-z]{1,49}")
}

func ValidPasswordLength(password string) bool {
  if len(password) < 6 {
    return false
  }

  return true
}
