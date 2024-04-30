package ratings

import (
	"errors"

	"github.com/94DanielBrown/roasts-api/internal/database"
)

func UpdateAverages(roastModels database.RoastModels, newReview database.Review) error {
	roast, err := roastModels.GetRoastByPrefix(newReview.RoastKey)
	if err != nil {
		return err
	}
	if roast == nil {
		return errors.New("no roast found")
	}

	newCount := roast.ReviewCount + 1
	roast.OverallRating = ((roast.OverallRating * float64(roast.ReviewCount)) + float64(newReview.OverallRating)) / float64(newCount)
	roast.MeatRating = ((roast.MeatRating * float64(roast.ReviewCount)) + float64(newReview.MeatRating)) / float64(newCount)
	roast.PotatoesRating = ((roast.PotatoesRating * float64(roast.ReviewCount)) + float64(newReview.PotatoesRating)) / float64(newCount)
	roast.VegRating = ((roast.VegRating * float64(roast.ReviewCount)) + float64(newReview.VegRating)) / float64(newCount)
	roast.GravyRating = ((roast.GravyRating * float64(roast.ReviewCount)) + float64(newReview.GravyRating)) / float64(newCount)
	// Repeat for all other ratings...

	roast.ReviewCount = newCount

	err = roastModels.UpdateRoast(roast)
	if err != nil {
		return err
	}

	return nil
}
