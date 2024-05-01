package rating

import (
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type (
	ChannelInfoFetcher interface {
		FetchMemeLink(memeID string, channelID string) string
		FetchMember(memberID MemberID) Member
	}

	SlackLinkFetcher struct {
		client *socketmode.Client
		logger zerolog.Logger
	}
)

func NewSlackLinkFetcher(client *socketmode.Client, logger zerolog.Logger) ChannelInfoFetcher {
	return SlackLinkFetcher{
		client: client,
		logger: logger,
	}
}

func (s SlackLinkFetcher) FetchMemeLink(memeID string, channelID string) string {
	permalink, err := s.client.GetPermalink(&slack.PermalinkParameters{
		Channel: channelID,
		Ts:      memeID,
	})

	if err != nil {
		s.logger.Err(err).
			Str("meme_id", memeID).
			Str("channel_id", channelID).
			Msg("⏳ Failed to fetch permalink to message")

		return ""
	}

	return permalink
}

func (s SlackLinkFetcher) FetchMember(memberID MemberID) Member {
	memberInfo, err := s.client.GetUserInfo(string(memberID))

	if err != nil {
		s.logger.Err(err).
			Str("member_id", string(memberID)).
			Msg("⏳ Failed to fetch member")

		return Member{}
	}

	return NewMember(MemberID(memberInfo.ID), memberInfo.Profile.RealName, memberInfo.Profile.DisplayName)
}
