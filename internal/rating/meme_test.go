package rating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaction_Score(t *testing.T) {
	t.Parallel()

	t.Run("When reaction is allowed, then score must be > 0", func(t *testing.T) {
		allowedReactions := []Reaction{"+1", "kekw", "omegalul"}

		for _, reaction := range allowedReactions {
			assert.True(t, reaction.Score() > 0)
		}
	})

	t.Run("When reaction is not allowed, then score must be equal to 0", func(t *testing.T) {
		allowedReactions := []Reaction{"-1", "kekl", "huh"}

		for _, reaction := range allowedReactions {
			assert.Equal(t, 0, reaction.Score())
		}
	})
}

func TestReactions_Score(t *testing.T) {
	t.Parallel()

	t.Run("When calculate over all reactions, then score every reaction", func(t *testing.T) {
		reactions := NewReactions([]Reaction{"+1", "kekw", "omegalul", "-1"})

		assert.Equal(t, 6, reactions.Score())
	})
}
