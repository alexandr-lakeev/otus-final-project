package internallogger

import (
	"bufio"
	"os"
	"path"
	"testing"

	"github.com/alexandr-lakeev/otus-final-project/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("Debug", func(t *testing.T) {
		tests := []struct {
			level      string
			messageKey int
		}{
			{level: "DEBUG", messageKey: 0},
			{level: "INFO", messageKey: 1},
			{level: "WARNING", messageKey: 2},
			{level: "ERROR", messageKey: 3},
			{level: "PANIC", messageKey: 4},
		}

		messages := []string{
			0: "debug msg",
			1: "info msg",
			2: "warning msg",
			3: "error msg",
			4: "panic msg",
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.level, func(t *testing.T) {
				stdout := path.Join(os.TempDir(), "/stdout")

				os.Stdout, _ = os.Create(stdout)

				logg, err := New(config.LoggerConf{
					Env:   "dev",
					Level: tc.level,
				})

				require.NoError(t, err)

				logg.Debug(messages[0])
				logg.Info(messages[1])
				logg.Warning(messages[2])
				logg.Error(messages[3])

				require.Panics(t, func() {
					logg.Panic(messages[4])
				})

				log, err := os.Open(stdout)
				require.NoError(t, err)

				scanner := bufio.NewScanner(log)
				n := tc.messageKey
				for scanner.Scan() {
					require.Contains(t, scanner.Text(), messages[n])
					n++
				}
			})
		}
	})

	t.Run("Wrong debug level", func(t *testing.T) {
		_, err := New(config.LoggerConf{
			Env:   "dev",
			Level: "WRONG_LEVEL",
		})

		require.Error(t, err)
	})
}
