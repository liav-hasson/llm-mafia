// This file containes pure voting functions

package domain

// count votes and return dict with player-id -> vote-count
func TallyVotes(votes map[string]string) map[string]int {
	// init empty map
	tally := make(map[string]int)

	for _, target := range votes {

		// check if key doesnt exist yet
		// note: could use 'tally[target]++' instead of the if statement
		count, exists := tally[target]
		if !exists {
			// add new target with 0 votes
			count = 0
		}

		// add vote to target
		tally[target] = count + 1
	}

	return tally
}

// getTopVoted is an internal helper that returns the player(s) with most votes
// unexported (lowercase) since it's only used internally
func getTopVoted(votes map[string]string) []string {
	// early exit if map is empty
	if len(votes) == 0 {
		return nil
	}

	tally := TallyVotes(votes)
	highestVotes := 0
	highestVotedPlayer := []string{}

	for player, votecount := range tally {

		// if found player with more votes, reset slice and append new
		if votecount > highestVotes {
			highestVotes = votecount

			// reset the slice and append the new player
			highestVotedPlayer = nil
			highestVotedPlayer = append(highestVotedPlayer, player)

			// if found player with same votes, append to slice
		} else if votecount == highestVotes {
			highestVotedPlayer = append(highestVotedPlayer, player)
		}
	}

	return highestVotedPlayer
}

// count votes and return the player with most votes
func GetVoteWinner(votes map[string]string) (string, bool) {
	highestVotedPlayer := getTopVoted(votes)

	// return empty string if no votes or tie
	if len(highestVotedPlayer) != 1 {
		return "", false
	}

	// return first value
	return highestVotedPlayer[0], true
}
