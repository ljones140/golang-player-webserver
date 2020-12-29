package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/ljones140/golang-player-webserver"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
	FinishCalled bool
}

func (g *GameSpy) Start(numberOfPlayers int, alertsDestination io.Writer) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
	g.FinishCalled = true
}

func TestCLI(t *testing.T) {
	t.Run("starts game with given numebr of players and records winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("1\nChris wins\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game.StartedWith, 1)
		assertFinishCalledWith(t, game.FinishedWith, "Chris")
	})

	t.Run("it does not start game when a non numeric value is entered", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Non Numeric\n")

		game := &GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it does not finish game if winner winner entered incorrectly", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		game := &GameSpy{}
		in := strings.NewReader("1\nChris Incorrectly Entered String\n")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessageSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})
}

func assertGameStartedWith(t *testing.T, got, numberOfPlayerswanted int) {
	t.Helper()
	if got != numberOfPlayerswanted {
		t.Errorf("expected game to be started with %d, but got %d", got, numberOfPlayerswanted)
	}
}

func assertFinishCalledWith(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected game to be finished with %q, but got %q", got, want)
	}
}

func assertGameNotStarted(t *testing.T, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("Should not have started game")
	}
}

func assertGameNotFinished(t *testing.T, game *GameSpy) {
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
