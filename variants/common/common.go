package common

import (
	"github.com/zond/godip"
	"github.com/zond/godip/orders"
	"github.com/zond/godip/state"
)

// Variant defines a dippy variant supported by godip.
type Variant struct {
	// Name is the display name and key for this variant.
	Name string
	// Start returns a state with the correct graph, units, phase and supply centers for starting this variant.
	Start func() (*state.State, error) `json:"-"`
	// BlankStart returns a state with the correct graph, phase and supply centers for starting this variant.
	BlankStart func() (*state.State, error) `json:"-"`
	// Blank returns a state with the correct graph and the provided phase for this variant.
	Blank func(godip.Phase) *state.State `json:"-"`
	// Phase returns a phase with the provided year, season and phase type for this variant.
	Phase func(int, godip.Season, godip.PhaseType) godip.Phase `json:"-"`
	// Parser for orders in the variant.
	Parser orders.Parser `json:"-"`
	// Graph is the graph for this variant.
	Graph func() godip.Graph `json:"-"`
	// Nations are the nations playing this variant.
	Nations []godip.Nation
	// PhaseTypes are the phase types the phases of this variant have.
	PhaseTypes []godip.PhaseType
	// Seasons are the seasons the phases of this variant have.
	Seasons []godip.Season
	// UnitTypes are the types the units of this variant have.
	UnitTypes []godip.UnitType
	// Function to return a nation with a solo (or the empty string if no such nation exists).
	SoloWinner func(*state.State) godip.Nation `json:"-"`
	// SVG representing the variant map graphics.
	SVGMap func() ([]byte, error) `json:"-"`
	// A version for the vector graphics (for use in caching mechanisms).
	SVGVersion string
	// SVG representing the variant units.
	SVGUnits map[godip.UnitType]func() ([]byte, error) `json:"-"`
	// Who the version was created by (or the empty string if no creator information is known).
	CreatedBy string
	// Version of the variant (or the empty string if no version information is known).
	Version string
	// A short description summarising the variant.
	Description string
	// The rules of the variant (in particular where they differ from classical).
	Rules string
}

// Return a function that declares a solo winner if a nation has more SCs than the given number (and more than any other nation).
func SCCountWinner(soloSupplyCenters int) func(*state.State) godip.Nation {
	return func(s *state.State) godip.Nation {
		// Create a map from nation to count of owned SCs.
		scCount := map[godip.Nation]int{}
		for _, nat := range s.SupplyCenters() {
			if nat != "" {
				scCount[nat] = scCount[nat] + 1
			}
		}
		// Check if there's a current clear leader.
		highestSCCount := 0
		var leader godip.Nation
		for nat, count := range scCount {
			if count > highestSCCount {
				leader = nat
				highestSCCount = count
			} else if count == highestSCCount {
				leader = ""
			}
		}
		// Return the nation if they have more than the required SCs.
		if leader != "" && scCount[leader] >= soloSupplyCenters {
			return leader
		}
		return ""
	}
}
