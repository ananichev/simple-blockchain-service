package service

import (
  "strconv"

  "github.com/ananichev/simple-blockchain-service/types"
)

func (s *Service) StoreLink(link types.Link) error {
  return s.storage.StoreLink(link)
}

func (s *Service) Status() (types.Status, error) {
  sts := types.Status{
    Id: s.myId,
    Name: s.myName,
    URL: s.myURL,
    LastHash: s.prevHash(),
    Neighbours: []string{},
  }

  links, err := s.storage.GetLinks()
  if err != nil {
    return types.Status{}, err
  }

  for _, l := range links {
    sts.Neighbours = append(sts.Neighbours, strconv.Itoa(l.Id))
  }
  return sts, nil
}
