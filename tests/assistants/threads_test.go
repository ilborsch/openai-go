package assistants

import (
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateThread_Happy(t *testing.T) {
	s := suite.New(t)

	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)
	require.NotEmpty(t, threadID)
}
