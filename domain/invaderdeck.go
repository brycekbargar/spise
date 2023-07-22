package domain

type InvaderDeck struct {
	game *Game

	Drawn  []InvaderCardDrawn
	InDeck []InvaderCardInDeck
}

func NewInvaderDeck(game *Game) *InvaderDeck {
	initial := []InvaderCard{
		StageOneUnknown,
		StageOneUnknown,
		StageOneUnknown,
		StageTwoUnknown,
		StageTwoUnknown,
		StageTwoUnknown,
		StageTwoUnknown,
		StageThreeUnknown,
		StageThreeUnknown,
		StageThreeUnknown,
		StageThreeUnknown,
		StageThreeUnknown,
	}

	mod, ok := modinvaderdec[game.SupportingAdversary]
	if ok {
		initial = mod(initial, game.SupportingAdversaryLevel)
	}
	mod, ok = modinvaderdec[game.LeadingAdversary]
	if ok {
		initial = mod(initial, game.LeadingAdversaryLevel)
	}

	deck := make([]InvaderCardInDeck, 0, len(initial))
	for _, c := range initial {
		deck = append(deck, InvaderCardInDeck{c})
	}

	return &InvaderDeck{
		game: game,

		Drawn:  []InvaderCardDrawn{},
		InDeck: deck,
	}
}

// InvaderCardInDeck wraps possibilites for the invader deck.
type InvaderCardInDeck struct {
	Card InvaderCard
}

// InvaderCardDrawn wraps drawn invader cards.
type InvaderCardDrawn struct {
	Card InvaderCard
}

var modinvaderdec = map[Adversary]func(deck []InvaderCard, lvl int) []InvaderCard{
	BrandenburgPrussia: func(deck []InvaderCard, lvl int) []InvaderCard {
		if lvl >= 2 {
			// Put 1 of the Stage III cards between Stage I and Stage II.
			s2ix, s3ix := -1, -1
			for i, c := range deck {
				if c.Stage == 2 && s2ix == -1 {
					s2ix = i
				}
				if c.Stage == 3 {
					s3ix = i
				}
			}

			if s2ix > 0 && s3ix > 0 && s2ix < s3ix {
				mod := make([]InvaderCard, len(deck[:s2ix]))
				copy(mod, deck[:s2ix])
				mod = append(mod, deck[s3ix])         // nozero
				mod = append(mod, deck[s2ix:s3ix]...) // nozero
				mod = append(mod, deck[s3ix+1:]...)   // nozero
				deck = mod
			}
		}

		if lvl >= 3 {
			// Remove an additional Stage I Card.
			for i, c := range deck {
				if c.Stage == 1 {
					deck = append(deck[:i], deck[i+1:]...)

					break
				}
			}
		}

		if lvl >= 4 {
			// Remove an additional Stage II Card.
			for i, c := range deck {
				if c.Stage == 2 {
					deck = append(deck[:i], deck[i+1:]...)

					break
				}
			}
		}

		if lvl >= 5 {
			// Remove an additional Stage I Card.
			for i, c := range deck {
				if c.Stage == 1 {
					deck = append(deck[:i], deck[i+1:]...)

					break
				}
			}
		}

		if lvl >= 6 {
			// Remove all Stage I Cards.
			for true {
				s1ix := -1
				for i, c := range deck {
					if c.Stage == 1 {
						s1ix = i

						break
					}
				}

				if s1ix != -1 {
					deck = append(deck[:s1ix], deck[s1ix+1:]...)

					continue
				}

				break
			}
		}

		return deck
	},
	HabsburgLivestock: func(deck []InvaderCard, lvl int) []InvaderCard {
		if lvl >= 3 {
			// Remove 1 additional Stage I Card.
			for i, c := range deck {
				if c.Stage == 1 {
					deck = append(deck[:i], deck[i+1:]...)

					break
				}
			}
		}

		return deck
	},
	HabsburgMines: func(deck []InvaderCard, lvl int) []InvaderCard {
		if lvl >= 4 {
			// Place the 'Salt Deposits' card in place of the 2nd Stage II card.
			s2 := 0
			for cix, c := range deck {
				if c.Stage == 2 {
					s2++
					if s2 == 2 {
						mod := make([]InvaderCard, len(deck[:cix]))
						copy(mod, deck[:cix])
						mod = append(mod, StageTwoSaltDeposits) // nozero
						mod = append(mod, deck[cix+1:]...)      // nozero
						deck = mod

						break
					}
				}
			}
		}

		return deck
	},
	Russia: func(deck []InvaderCard, lvl int) []InvaderCard {
		if lvl >= 4 {
			// Put 1 Stage III Card after each Stage II Card.
			os2ix := -1
			for true {
				s3ix := -1
				for i, c := range deck {
					if c.Stage == 3 {
						s3ix = i
					}
				}

				s2ix := -1
				for i, c := range deck {
					if c.Stage == 2 && i > os2ix {
						s2ix = i

						break
					}
				}

				if s2ix != -1 && s3ix != -1 {
					os2ix = s2ix
					mod := make([]InvaderCard, len(deck[:s2ix+1]))
					copy(mod, deck[:s2ix+1])
					mod = append(mod, deck[s3ix])           // nozero
					mod = append(mod, deck[s2ix+1:s3ix]...) // nozero
					mod = append(mod, deck[s3ix+1:]...)     // nozero
					deck = mod

					continue
				}

				break
			}
		}

		return deck
	},
	Scotland: func(deck []InvaderCard, lvl int) []InvaderCard {
		if lvl >= 2 {
			// Place "Coastal Lands" as the 3rd Stage II card.
			stage2 := 0
			for cix, c := range deck {
				if c.Stage == 2 {
					stage2++
					if stage2 == 3 {
						mod := make([]InvaderCard, len(deck[:cix]))
						copy(mod, deck[:cix])
						mod = append(mod, StageTwoCoastal) // nozero
						mod = append(mod, deck[cix+1:]...) // nozero
						deck = mod

						break
					}
				}
			}

			// Move the two Stage II Cards above it up by one.
			stage2 = 0
			s2ix1, s2ix2 := -1, -1
			for cix, card := range deck {
				card := card
				if card == StageTwoCoastal {
					break
				}

				if card.Stage == 2 {
					stage2++
					switch stage2 {
					case 1:
						s2ix1 = cix
					case 2:
						s2ix2 = cix
					}
					if stage2 >= 2 {
						break
					}
				}
			}

			if s2ix1 >= 1 {
				mod := make([]InvaderCard, len(deck[:s2ix1-1]))
				copy(mod, deck[:s2ix1-1])
				mod = append(mod, deck[s2ix1])       // nozero
				mod = append(mod, deck[s2ix1-1])     // nozero
				mod = append(mod, deck[s2ix1+1:]...) // nozero
				deck = mod
			}
			if s2ix2 >= 2 {
				mod := make([]InvaderCard, len(deck[:s2ix2-1]))
				copy(mod, deck[:s2ix2-1])
				mod = append(mod, deck[s2ix2])       // nozero
				mod = append(mod, deck[s2ix2-1])     // nozero
				mod = append(mod, deck[s2ix2+1:]...) // nozero
				deck = mod
			}
		}

		if lvl >= 4 {
			// Replace the bottom Stage I Card with the bottom Stage III
			s1ix, s3ix := -1, -1
			for cix, c := range deck {
				if c.Stage == 1 {
					s1ix = cix
				}
				if c.Stage == 3 {
					s3ix = cix
				}
			}

			if s1ix >= 0 && s3ix >= 1 {
				mod := make([]InvaderCard, len(deck[:s1ix]))
				copy(mod, deck[:s1ix])
				mod = append(mod, deck[s3ix])           // nozero
				mod = append(mod, deck[s1ix+1:s3ix]...) // nozero
				mod = append(mod, deck[s3ix+1:]...)     // nozero
				deck = mod
			}
		}

		return deck
	},
}
