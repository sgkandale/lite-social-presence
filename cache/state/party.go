package state

import "context"

func (c *Client) PutPartyMembersList(ctx context.Context, partyName string, members []string) error {
	c.cache.Set(
		PatyMembersKey(partyName),
		members,
		1,
	)
	return nil
}

func (c *Client) GetPartyMembersList(ctx context.Context, partyName string) ([]string, error) {
	membersList, found := c.cache.Get(PatyMembersKey(partyName))
	if found || membersList != nil {
		return membersList.([]string), nil
	}
	return nil, nil
}
