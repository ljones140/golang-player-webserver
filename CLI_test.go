package poker_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	poker "github.com/ljones140/golang-player-webserver"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("starts game with given numebr of players and records winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("1\nChris wins\n")
		game := &poker.GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 1)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("it does not start game when a non numeric value is entered", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Non Numeric\n")

		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it does not finish game if winner winner entered incorrectly", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		game := &poker.GameSpy{}
		in := strings.NewReader("1\nChris Incorrectly Entered String\n")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})
}

func assertGameStartedWith(t *testing.T, game *poker.GameSpy, numberOfPlayerswanted int) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == numberOfPlayerswanted
	})

	if !passed {
		t.Errorf("expected game to be started with %d, but got %d", game.StartedWith, numberOfPlayerswanted)
	}
}

func assertFinishCalledWith(t *testing.T, game *poker.GameSpy, winner string) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == winner
	})

	if !passed {
		t.Errorf("expected game to be finished with %q, but got %q", game.FinishedWith, winner)
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func assertGameNotStarted(t *testing.T, game *poker.GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("Should not have started game")
	}
}

func assertGameNotFinished(t *testing.T, game *poker.GameSpy) {
	t.Helper()
	if game.FinishCalled {
		t.Errorf("Should not have started game")
	}
}
func assertMessageSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()

	if got != want {
		t.Errorf("got %q sent to stdout, wanted %q", got, messages)
	}
}
