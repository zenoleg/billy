package rating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaction_Score(t *testing.T) {
	t.Parallel()

	t.Run("When reaction is allowed, then score must be > 0", func(t *testing.T) {
		testCases := []struct {
			reaction      string
			count         int
			expectedScore int
		}{
			{reaction: "+1", count: 10, expectedScore: 10},
			{reaction: "omegalul", count: 10, expectedScore: 20},
			{reaction: "kekw", count: 10, expectedScore: 30},
		}

		for _, testCase := range testCases {
			assert.Equal(t, testCase.expectedScore, NewReaction(testCase.reaction, testCase.count).Score())
		}
	})

	t.Run("When reaction is not allowed, then score must be equal to 0", func(t *testing.T) {
		allowedReactions := []string{"-1", "kekl", "huh"}

		for _, reaction := range allowedReactions {
			assert.Equal(t, 0, NewReaction(reaction, 1).Score())
		}
	})
}

func TestReactions_Score(t *testing.T) {
	t.Parallel()

	t.Run("When calculate over all reactions, then score every reaction", func(t *testing.T) {
		reactions := NewReactions([]Reaction{
			{reaction: "+1", count: 1},
			{reaction: "kekw", count: 3},
			{reaction: "shit", count: 33},
		})

		assert.Equal(t, 10, reactions.Score())
	})
}
