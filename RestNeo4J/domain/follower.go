package domain

import (
	"encoding/json"
	"io"
)

type UserFollower struct {
	UserId            int32 `json:"userId,omitempty"`
	Username          string `json:"username,omitempty"`
	FollowingUserId   int32 `json:"followingUserId,omitempty"`
	FollowingUsername string `json:"followingUsername,omitempty"`
}

func (o *UserFollower) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *UserFollower) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
