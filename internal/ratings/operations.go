package ratings

import "github.com/94DanielBrown/roasts/internal/database"

func UpdateAverages(roastModels database.RoastModels, roastID string, newReview database.Review) error {
	// Step 1: Fetch the existing roast item
	roast, err := roastModels.GetRoastByPrefix(roastID)
	if err != nil {
		return err
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
