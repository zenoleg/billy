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
			{reaction: "omegalul", count: 10, expectedScore: 10},
			{reaction: "kekw", count: 10, expectedScore: 10},
		}

		for _, testCase := range testCases {
			assert.Equal(t, testCase.expectedScore, NewReaction(testCase.reaction, testCase.count).Score())
		}
	})
}

func TestReactions_Score(t *testing.T) {
	t.Parallel()

	t.Run("When calculate over all reactions, then score every reaction", func(t *testing.T) {
		reactions := NewReactions([]Reaction{
			{reaction: "+1", count: 1},
			{reaction: "kekw", count: 1},
			{reaction: "shit", count: 1},
		})

		assert.Equal(t, 3, reactions.Score())
	})
}
