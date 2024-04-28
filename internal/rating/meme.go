package rating

import "slices"

type (
	MemberID  string
	PostID    string
	Reaction  string
	Reactions []Reaction

	Member struct {
		id   MemberID
		name string
	}

	Post struct {
		postID string
		from   MemberID
		score  int
	}
)

func NewPost(postID string, from MemberID, reactions Reactions) Post {
	return Post{postID: postID, from: from, score: reactions.Score()}
}

func NewMember(ID MemberID, name string) Member {
	return Member{id: ID, name: name}
}

func NewReaction(value string) Reaction {
	return Reaction(value)
}

func NewReactions(reactions []Reaction) Reactions {
	return Reactions(reactions)
}

func (r Reaction) Score() int {
	if slices.Contains([]Reaction{"+1"}, r) {
		return 1
	}
	if slices.Contains([]Reaction{"omegalul"}, r) {
		return 2
	}
	if slices.Contains([]Reaction{"kekw"}, r) {
		return 3
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
