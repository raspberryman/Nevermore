package cmd

import (
	"strconv"
)

// Syntax: DROP item
func init() {
	addHandler(drop{},"DROP")
	addHelp("Usage:  drop itemName # \n \n Drop the specified item name and number.", 0, "drop")
}

type drop cmd

func (drop) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to drop?")
		return
	}

	// We have at least 2 items here so lets move forward with that
	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) < 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}else{
			s.msg.Actor.SendInfo("What did you want to drop?")
			return
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to drop?")
		return
	}

	where := s.where.Items

	s.actor.Inventory.Lock()
	where.Lock()
	s.actor.Inventory.Remove(target)
	where.Add(target)
	s.actor.Inventory.Unlock()
	where.Unlock()


	s.msg.Actor.SendGood("You drop ", target.Name, ".")
	s.msg.Observer.SendInfo( s.actor.Name, " drops ", target.Name, ".")

	s.ok = true
}
