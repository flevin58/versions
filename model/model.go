package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/flevin58/versions/cfg"
)

type InstalledStatus int

const (
	homebrew = iota // Macos installer
	scoop           // Windows installer
	other
	unavailable
)

type cmdData struct {
	command   string
	version   string
	where     string
	installed InstalledStatus
}

func (c *cmdData) Record() []string {
	return []string{c.command, c.version, c.where}
}

var (
	cmdHeader cmdData = cmdData{"Command", "Version", "Where", unavailable}
	cmdList   []cmdData
)

func getVersion(cmd cfg.Command) string {
	const notavailable = "not available"
	command := exec.Command(cmd.Name, cmd.VersionFlag)
	out, err := command.Output()
	if err != nil {
		return notavailable
	}
	version := strings.Split(string(out), "\n")[cmd.VersionLine-1]
	re := regexp.MustCompile(`^[^\d]*([0-9.]+)`)
	match := re.FindStringSubmatch(version)
	result := match[1]
	return result
}

func Add(cmd cfg.Command) {
	var data cmdData
	where, err := exec.LookPath(cmd.Name)
	if err != nil {
		data = cmdData{
			command:   cmd.Name,
			version:   "",
			where:     "",
			installed: unavailable,
		}
	} else {
		var inst InstalledStatus = other
		switch runtime.GOOS {
		case "windows":
			winPrefix := os.ExpandEnv("$USERPROFILE\\scoop")
			fmt.Println(where)
			fmt.Println(winPrefix)
			if strings.HasPrefix(where, winPrefix) {
				inst = scoop
			}
		case "darwin":
			macPrefix := "/opt/homebrew"
			if strings.HasPrefix(where, macPrefix) {
				inst = homebrew
			}
		}
		data = cmdData{
			command:   cmd.Name,
			version:   getVersion(cmd),
			where:     where,
			installed: inst,
		}
	}
	cmdList = append(cmdList, data)
}

func ToCSV(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	header := cmdHeader.Record()
	if err = w.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}
	for i := range cmdList {
		record := cmdList[i].Record()
		if err = w.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}
	return nil
}

func ToTable() {
	// Prepare the table
	HeaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).AlignHorizontal(lipgloss.Center)
	EvenRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fff200"))
	OddRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	table := table.New().Border(lipgloss.NormalBorder())
	table.Headers("Command", "Version", "Where")
	table.StyleFunc(func(row, col int) lipgloss.Style {
		switch {
		case row == 0:
			return HeaderStyle
		case row%2 == 0:
			return EvenRowStyle
		default:
			return OddRowStyle
		}
	})

	// We assume that as "nerds" we are installing using the following:
	// macos: homebrew
	// windows: scoop
	var emoji string
	for _, row := range cmdList {
		switch row.installed {
		case homebrew:
			emoji = "🍺 "
		case scoop:
			emoji = "🍨 "
		case other:
			emoji = "🟢 "
		default:
			emoji = "🔴 "
		}

		table.Row(emoji+row.command, row.version, row.where)
	}
	fmt.Println(table.Render())
}

func ToText() {
}
