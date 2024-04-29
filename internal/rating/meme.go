package rating

import "slices"

type (
	MemberID string
	Reaction struct {
		reaction string
		count    int
	}
	Reactions []Reaction

	Member struct {
		id   MemberID
		name string
	}

	Meme struct {
		id        string
		channelID string
		memberID  MemberID
		score     int
		timestamp string
	}
)

func NewMeme(id string, channelID string, from MemberID, reactions Reactions, timestamp string) Meme {
	return Meme{
		id:        id,
		channelID: channelID,
		memberID:  from,
		score:     reactions.Score(),
		timestamp: timestamp,
	}
}

func (m Meme) Rate(score int) Meme {
	m.score = m.score + score

	return m
}

func NewMemberID(value string) MemberID {
	return MemberID(value)
}

func NewMember(ID MemberID, name string) Member {
	return Member{id: ID, name: name}
}

func NewReaction(reaction string, count int) Reaction {
	return Reaction{
		reaction: reaction,
		count:    count,
	}
}

func NewReactions(reactions []Reaction) Reactions {
	return Reactions(reactions)
}

func (r Reaction) Score() int {
	if slices.Contains([]string{"+1"}, r.reaction) {
		return 1 * r.count
	}
	if slices.Contains([]string{"omegalul"}, r.reaction) {
		return 2 * r.count
	}
	if slices.Contains([]string{"kekw"}, r.reaction) {
		return 3 * r.count
	}

	return 0
}

func (r Reactions) Score() int {
	score := 0

	for _, reaction := range r {
		score = score + reaction.Score()
	}

	return score
}
