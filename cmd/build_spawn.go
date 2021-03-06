package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
	"strings"
)

func init() {
	addHandler(spawn{}, "spawn")
	addHelp("Usage:  spawn (mob|item) (name) \n \n Use this command to spawn a mob or item to be modified: \n" +
		"Items: Item will be added to your inventory\n" +
		"  -->  If you wish to save it as the template for that item, use the 'savetemplate item' command\n" +
		"Mob:  Mob will be spawned into your room. \n" +
		"  -->  If you wish to save it as the template for that mob, use the 'savetemplate mob' command\n\n",50, "spawn")
}

type spawn cmd

func (spawn) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 50 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Spawn what?")
		return
	}

	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "mob":
		//log.Println("Trying to do a spawn...")
		mob_id, err := strconv.Atoi(s.words[1])
		if err != nil {
			s.msg.Actor.SendBad("What mob ID do you want to spawn?")
			return
		}
		//log.Println("Copying mob")
		newMob := objects.Mobs[int64(mob_id)]
		//#log.Println("Adding the mob to the room...")
		s.msg.Actor.Send(text.Magenta + "You encounter: " + newMob.Name + text.Reset )
		s.msg.Observer.Send(text.Magenta + "You encounter: " + newMob.Name + text.Reset )
		s.where.Mobs.Add(newMob)
	// Handle Exits
	case "object":
		return
	default:
		s.msg.Actor.SendBad("Not an object that can be edited, or WIP")
	}

	s.ok = true
	return
}
