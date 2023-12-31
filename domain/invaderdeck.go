package domain

import (
	"errors"
)

// ErrNoInvaderCard occurs when the an action tries to draw from an empty invader deck.
var (
	ErrNoInvaderCard = errors.New("the invader deck is empty")
	// ErrInvaderCardNotReturnable occurs when Fractured Days returns an invalid card.
	ErrInvaderCardNotReturnable = errors.New(
		"the invader card must be within one stage to return",
	)
	// ErrNotEntrenched occurs when Russia's game effect is active.
	ErrNotEntrenched = errors.New(
		"only Russia 5+ is entrenched in the face of fear",
	)
)

type InvaderDeck struct {
	game *Game

	Drawn  []InvaderCardDrawn
	InDeck []InvaderCardInDeck
}

func NewInvaderDeck(game *Game) *InvaderDeck {
	initial := []InvaderCardInDeck{
		{StageOneUnknown, false},
		{StageOneUnknown, false},
		{StageOneUnknown, false},
		{StageTwoUnknown, false},
		{StageTwoUnknown, false},
		{StageTwoUnknown, false},
		{StageTwoUnknown, false},
		{StageThreeUnknown, false},
		{StageThreeUnknown, false},
		{StageThreeUnknown, false},
		{StageThreeUnknown, false},
		{StageThreeUnknown, false},
	}

	mod, ok := modinvaderdec[game.SupportingAdversary]
	if ok {
		initial = mod(initial, game.SupportingAdversaryLevel)
	}
	mod, ok = modinvaderdec[game.LeadingAdversary]
	if ok {
		initial = mod(initial, game.LeadingAdversaryLevel)
	}

	return &InvaderDeck{
		game: game,

		Drawn:  []InvaderCardDrawn{},
		InDeck: initial,
	}
}

func (deck *InvaderDeck) Draw(card InvaderCard) error {
	if len(deck.InDeck) == 0 {
		return ErrNoInvaderCard
	}

	if card.Stage != deck.InDeck[0].Stage {
		return ErrInvalidInvaderCard
	}

	deck.Drawn = append(deck.Drawn, InvaderCardDrawn{card, false})
	deck.InDeck = deck.InDeck[1:]
	deck.setReturnable()

	return nil
}

func (deck *InvaderDeck) Return(card InvaderCard) error {
	for cix, c := range deck.Drawn {
		if c.InvaderCard == card {
			if !c.PastReturnable {
				return ErrInvaderCardNotReturnable
			}

			mod := make([]InvaderCardDrawn, len(deck.Drawn[:cix]))
			copy(mod, deck.Drawn[:cix])
			mod = append( // nozero
				mod,
				InvaderCardDrawn{deck.InDeck[0].InvaderCard, true},
			)
			mod = append( // nozero
				mod,
				deck.Drawn[cix+1:]...)
			deck.Drawn = mod

			deck.InDeck[0] = InvaderCardInDeck{card, false}
			deck.setReturnable()

			return nil
		}
	}

	return ErrInvalidInvaderCard
}

func (deck *InvaderDeck) IgnoreRisingInterest() {
	if len(deck.InDeck) == 0 ||
		deck.InDeck[0].SpeciallyPlaced {
		return
	}

	deck.InDeck = deck.InDeck[1:]
	deck.setReturnable()
}

func (deck *InvaderDeck) DistractHardworkingSettlers() {
	s2ix, s3ix := -1, -1
	for i := len(deck.InDeck) - 1; i >= 0; i-- {
		if s2ix != -1 && deck.InDeck[i].Stage == 2 {
			s2ix = i
		}
		if s3ix != -1 && deck.InDeck[i].Stage == 3 {
			s3ix = i
		}
	}

	if s2ix != -1 && s3ix != -1 && s2ix > s3ix {
		// This is kind of a hack.
		s2ix, s3ix = s3ix, s2ix
	}

	if s3ix != -1 {
		deck.InDeck = append(deck.InDeck[:s3ix], deck.InDeck[s3ix+1:]...)
	}
	if s2ix != -1 {
		deck.InDeck = append(deck.InDeck[:s2ix], deck.InDeck[s2ix+1:]...)
	}
}

func (deck *InvaderDeck) Entrenched(card InvaderCard) error {
	if (deck.game.LeadingAdversary != Russia || deck.game.LeadingAdversaryLevel < 5) &&
		(deck.game.SupportingAdversary != Russia || deck.game.SupportingAdversaryLevel < 5) {
		return ErrNotEntrenched
	}

	if card.Stage != 2 && card.Stage != 3 {
		return ErrInvalidInvaderCard
	}

	deck.Drawn = append(deck.Drawn, InvaderCardDrawn{card, false})
	deck.setReturnable()

	return nil
}

// InvaderCardInDeck wraps possibilites for the invader deck.
type InvaderCardInDeck struct {
	InvaderCard
	SpeciallyPlaced bool
}

// InvaderCardDrawn wraps drawn invader cards.
type InvaderCardDrawn struct {
	InvaderCard
	PastReturnable bool
}

func (deck *InvaderDeck) setReturnable() {
	rem := len(deck.InDeck) != 0
	stg := -99
	if rem {
		stg = deck.InDeck[0].Stage
	}
	for i := range deck.Drawn {
		c := deck.Drawn[i]
		if c.Stage > stg {
			deck.Drawn[i].PastReturnable = rem && (c.Stage-stg <= 1)
		} else {
			deck.Drawn[i].PastReturnable = rem && (stg-c.Stage <= 1)
		}
	}
}

// TODO: what a mess.
var modinvaderdec = map[Adversary]func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck{
	BrandenburgPrussia: func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck {
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
				mod := make([]InvaderCardInDeck, len(deck[:s2ix]))
				copy(mod, deck[:s2ix])
				deck[s3ix].SpeciallyPlaced = true
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
	HabsburgLivestock: func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck {
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
	HabsburgMines: func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck {
		if lvl >= 4 {
			// Place the 'Salt Deposits' card in place of the 2nd Stage II card.
			s2 := 0
			for cix, c := range deck {
				if c.Stage == 2 {
					s2++
					if s2 == 2 {
						mod := make([]InvaderCardInDeck, len(deck[:cix]))
						copy(mod, deck[:cix])
						mod = append( // nozero
							mod,
							InvaderCardInDeck{StageTwoSaltDeposits, true},
						)
						mod = append( // nozero
							mod,
							deck[cix+1:]...)
						deck = mod

						break
					}
				}
			}
		}

		return deck
	},
	Russia: func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck {
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
					mod := make([]InvaderCardInDeck, len(deck[:s2ix+1]))
					copy(mod, deck[:s2ix+1])
					deck[s3ix].SpeciallyPlaced = true
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
	Scotland: func(deck []InvaderCardInDeck, lvl int) []InvaderCardInDeck {
		if lvl >= 2 {
			// Place "Coastal Lands" as the 3rd Stage II card.
			stage2 := 0
			for cix, c := range deck {
				if c.Stage == 2 {
					stage2++
					if stage2 == 3 {
						mod := make([]InvaderCardInDeck, len(deck[:cix]))
						copy(mod, deck[:cix])
						mod = append( // nozero
							mod,
							InvaderCardInDeck{StageTwoCoastal, true},
						)
						mod = append( // nozero
							mod,
							deck[cix+1:]...)
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
				if card.InvaderCard == StageTwoCoastal {
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
				mod := make([]InvaderCardInDeck, len(deck[:s2ix1-1]))
				copy(mod, deck[:s2ix1-1])
				deck[s2ix1].SpeciallyPlaced = true
				mod = append(mod, deck[s2ix1])       // nozero
				mod = append(mod, deck[s2ix1-1])     // nozero
				mod = append(mod, deck[s2ix1+1:]...) // nozero
				deck = mod
			}
			if s2ix2 >= 2 {
				mod := make([]InvaderCardInDeck, len(deck[:s2ix2-1]))
				copy(mod, deck[:s2ix2-1])
				deck[s2ix2].SpeciallyPlaced = true
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
				mod := make([]InvaderCardInDeck, len(deck[:s1ix]))
				copy(mod, deck[:s1ix])
				deck[s3ix].SpeciallyPlaced = true
				mod = append(mod, deck[s3ix])           // nozero
				mod = append(mod, deck[s1ix+1:s3ix]...) // nozero
				mod = append(mod, deck[s3ix+1:]...)     // nozero
				deck = mod
			}
		}

		return deck
	},
}
