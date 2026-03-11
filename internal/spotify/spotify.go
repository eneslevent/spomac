package spotify

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type PlayerState struct {
	IsPlaying bool
	Track     string
	Artist    string
	Album     string
	Position  float64 // seconds
	Duration  float64 // seconds
}

func tell(command string) (string, error) {
	script := fmt.Sprintf(`tell application "Spotify" to %s`, command)
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func IsRunning() bool {
	script := `tell application "System Events" to (name of processes) contains "Spotify"`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func Launch() error {
	script := `tell application "Spotify" to activate`
	_, err := exec.Command("osascript", "-e", script).Output()
	return err
}

func GetState() (PlayerState, error) {
	if !IsRunning() {
		return PlayerState{}, fmt.Errorf("Spotify is not running")
	}

	script := `
tell application "Spotify"
	set pState to player state as string
	set tName to name of current track
	set tArtist to artist of current track
	set tAlbum to album of current track
	set pPos to player position
	set tDur to duration of current track
	return pState & "|||" & tName & "|||" & tArtist & "|||" & tAlbum & "|||" & (pPos as string) & "|||" & (tDur as string)
end tell`

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return PlayerState{}, fmt.Errorf("failed to get Spotify state: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(out)), "|||")
	if len(parts) != 6 {
		return PlayerState{}, fmt.Errorf("unexpected output: %s", string(out))
	}

	// macOS locale may use comma as decimal separator, replace with dot
	posStr := strings.Replace(parts[4], ",", ".", 1)
	durStr := strings.Replace(parts[5], ",", ".", 1)

	pos, _ := strconv.ParseFloat(posStr, 64)
	dur, _ := strconv.ParseFloat(durStr, 64)

	// Spotify returns duration in milliseconds
	durSec := dur / 1000.0

	return PlayerState{
		IsPlaying: parts[0] == "playing",
		Track:     parts[1],
		Artist:    parts[2],
		Album:     parts[3],
		Position:  pos,
		Duration:  durSec,
	}, nil
}

func PlayPause() error {
	_, err := tell("playpause")
	return err
}

func Next() error {
	_, err := tell("next track")
	return err
}

func Previous() error {
	_, err := tell("previous track")
	return err
}

func GetVolume() (int, error) {
	out, err := tell("sound volume")
	if err != nil {
		return 0, err
	}
	vol, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, err
	}
	return vol, nil
}

func SetVolume(vol int) error {
	if vol < 0 {
		vol = 0
	}
	if vol > 100 {
		vol = 100
	}
	_, err := tell(fmt.Sprintf("set sound volume to %d", vol))
	return err
}

