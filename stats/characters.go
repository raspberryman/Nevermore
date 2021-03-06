// Stats and global listing of characters.

package stats

import (
	"fmt"
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/text"
	"io"
	"log"
	"strconv"
	"sync"
)

// Currently active characters
type characterStats struct {
	sync.Mutex
	list []*objects.Character
}

var ActiveCharacters = &characterStats{}
var IpMap = make(map[string]string)

// Add adds the specified character to the list of characters.
func (c *characterStats) Add(character *objects.Character, address string) {
	if character.Flags["invisible"] || character.Class >=50 {
		c.MessageGM("###: " + character.Name + "["+ address +"] joins the realm.")
	}else{
		c.MessageAll("###: " + character.Name + " joins the realm.")
	}
	c.Lock()
	c.list = append(c.list, character)
	IpMap[character.Name] = address
	c.Unlock()
}

// Pass character as a pointer, compare and remove
func (c *characterStats) Remove(character *objects.Character) {
	log.Println("trying to let everyone know...")
	if character.Flags["invisible"] || character.Class >=50 {
		c.MessageGM("###: " + character.Name + " departs the realm.")
	}else{
		c.MessageAll("###: " + character.Name + " departs the realm.")
	}

	c.Lock()


	for i, p := range c.list {
		if p == character {
			copy(c.list[i:], c.list[i+1:])
			c.list[len(c.list)-1] = nil
			c.list = c.list[:len(c.list)-1]
			delete(IpMap, character.Name)
			break
		}
	}


	if len(c.list) == 0 {
		c.list = make([]*objects.Character, 0, 10)
	}

	c.Unlock()
}

func (c *characterStats) Find(name string) *objects.Character {
	c.Lock()
	for _, p := range c.list {
		if p.Name == name {
			c.Unlock()
			return p
		}
	}
	return nil
}

// List returns the names of all characters in the character list. The omit parameter
// may be used to specify a character that should be omitted from the list.
func (c *characterStats)  List() []string {
	c.Lock()

	list := make([]string, 0, len(c.list))


	for _, character := range c.list {
		if character.Flags["invisible"] == true {
			continue
		}

		if character.Title != "" {
			list = append(list, fmt.Sprintf("%s, %s, %s", character.Name, character.ClassTitle, character.Title))
		}else{
			list = append(list, fmt.Sprintf("%s, %s", character.Name, character.ClassTitle))
		}
	}

	c.Unlock()
	return list
}

// List returns the names of all characters in the character list. The omit parameter
// may be used to specify a character that should be omitted from the list.
func (c *characterStats)  GMList() []string {
	c.Lock()

	list := make([]string, 0, len(c.list))


	for _, character := range c.list {
		if character.Title != "" {
			list = append(list, fmt.Sprintf("(Room: %s) (%s) %s, %s, %s", strconv.Itoa(int(character.ParentId)), IpMap[character.Name], character.Name, character.ClassTitle, character.Title))
		}else{
			list = append(list, fmt.Sprintf("(Room: %s) (%s) %s, %s", strconv.Itoa(int(character.ParentId)), IpMap[character.Name], character.Name, character.ClassTitle))
		}
	}

	c.Unlock()
	return list
}

func (c *characterStats) MessageAll(msg string) {
	c.Lock()

	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, msg)
	players := []io.Writer{}
	for _, p := range c.list {
		players = append(players, p)
	}
	msgbuf.Deliver(players...)

	c.Unlock()
	return
}

func (c *characterStats) MessageGM(msg string) {
	c.Lock()

	// Setup buffer
	msgbuf := message.AcquireBuffer()
	msgbuf.Send(text.White, "[DM] ", msg)
	players := []io.Writer{}
	for _, p := range c.list {
		if p.Class >= int64(50) {
			players = append(players, p)
		}
	}
	msgbuf.Deliver(players...)

	c.Unlock()
	return
}


// Len returns the length of the character list.
func (c *characterStats) Len() (l int) {
	c.Lock()
	l = len(c.list)
	c.Unlock()
	return
}
