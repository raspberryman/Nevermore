// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package frontend

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"sort"
	"strings"
)

// menu embeds a frontend instance adding fields and methods specific to
// the main menu.
type start struct {
	*frontend
	optionEnd int
	powerCharacter string
	characters []string
	pwChange [16]byte
}

// NewMenu returns a menu with the specified frontend embedded. The returned
// menu can be used for processing the main menu and it's options.
func NewStart(f *frontend) (m *start){
	m = &start{frontend: f}
	m.startDisplay()
	return
}

// menuDisplay shows the main menu of options available once a player is logged
// into the system.
func (m *start) startDisplay() {
	// Load Characters
	m.characters, _ = data.ListChars(m.account)
	sort.Strings(m.characters)
	m.powerCharacter, _ = data.ListPowerChar(m.account)
	var output strings.Builder
	m.optionEnd = 2
	output.WriteString(text.White + `
 Message of the Day: 
` + config.Server.Motd + text.Good +`
 
=========
 Choose an action:
 ---------
 0. Quit
 X. Make a new character
 2. Change account password
`)
	switch m.accountType {
	case 100:
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a gamemaster account.\r\n")
		}else{
			output.WriteString(" 3. Load gamemaster account:" + m.powerCharacter + "\r\n")
		}
	case 50:
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a builder account.\r\n")
		}else{
			output.WriteString(" 3. Load builder account:" + m.powerCharacter + "\r\n")
		}
		m.optionEnd = 3
	case 60:
		if m.powerCharacter == "" {
			output.WriteString(" 3. Create a dungeon master account.\r\n")
		}else{
			output.WriteString(" 3. Load dungeonmaster account:" + m.powerCharacter + "\r\n")
		}
		m.optionEnd = 3
	default:
		m.optionEnd = 2
	}
	output.WriteString("\n==== Your Character List ====\n")
	output.WriteString("(Enter the name of the character you wish to play)\n")
	output.WriteString(m.characterList())

	m.buf.Send(output.String())
	m.nextFunc = m.startProcess
}

// menuProcess takes the current input and processes it as a menu option. If
// the option is valid the corresponding action is taken. If the option is not
// valid the player is notified and we wait for another option to be chosen.
func (m *start) startProcess() {
	switch string(m.input) {
	case "":
		return
	case "1":
		m.buf.Send(text.Bad, "New character creation is disabled at this time.", text.Reset)
		//NewCharacter(m.frontend)
	case "0":
		m.Close()
	case "2":
		m.buf.Send(text.Info, "Enter your new password:")
		m.nextFunc = m.verifyPw
	default:
		if m.accountType > 0 && string(m.input) == "3"{
			if m.powerCharacter == ""{
				NewPCharacter(m.frontend)
				return
			}else{
				StartGame(m.frontend, m.powerCharacter)
			}
		}else if utils.StringInLower(string(m.input), m.characters){
			StartGame(m.frontend, string(m.input))
		}
		m.buf.Send(text.Bad, "Invalid option selected. Please try again.", text.Reset)
	}
}

func (m *start) verifyPw(){
	switch l := len(m.input); {
	case l == 0:
		m.buf.Send(text.Bad, "No text sent, returning to menu.", text.Reset)
		m.startDisplay()
	default:
		m.pwChange = md5.Sum(m.input)
		m.buf.Send(text.Good, "Please type your password again to change.")
		m.nextFunc = m.changePw
	}
	return
}

func (m *start) changePw(){
	switch l := len(m.input); {
	case l == 0:
		m.buf.Send(text.Info, "Password change cancelled.\n", text.Reset)
		m.startDisplay()
	default:
		if md5.Sum(m.input) != m.pwChange {
			m.buf.Send(text.Bad, "Passwords do not match, please try again. Or enter nothing to cancel\n", text.Reset)
			m.nextFunc = m.changePw
		}else{
			data.UpdatePassword(m.account,hex.EncodeToString(m.pwChange[:]))
			m.buf.Send(text.Good, "Password changed!\n")
			m.startDisplay()
		}

	}
	return
}

func (m *start) characterList() string {
	var charList strings.Builder
	for _, charn := range m.characters {
		charList.Write([]byte(" " + charn))
	}
	/*
		// Width of gutter between columns
	const gutter = 4

	// Find longest key extracted
	maxWidth := 0
	for _, cmd:= range m.characters {
		if l := len(cmd)+4; l > maxWidth {
			maxWidth = l
		}
	}

	var (
		columnWidth = maxWidth + gutter
		columnCount = 80 / columnWidth
		rowCount    = len(m.characters) / columnCount
	)

	// If we have a partial row we need to account for it
	if len(m.characters) > rowCount*columnCount {
		rowCount++
	}

	// NOTE: cell is not (row * columnCount) + column as we are pivoting the
	// table so that the commands are alphabetical DOWN the rows not across the
	// columns.
	for row := 0; row < rowCount; row++ {
		line := []byte{}
		for column := 0; column < columnCount; column++ {
			cell := (column * rowCount) + row
			if cell < len(m.characters) {
				line = append(line, m.characters[cell]...)
				line = append(line, strings.Repeat(" ", columnWidth-len(m.characters[cell]))...)
			}
		}
		charList.Write(line)
	}
	*/
	return charList.String()

}