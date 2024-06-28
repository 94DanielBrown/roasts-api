package ratings

import (
	"errors"
	"fmt"

	"github.com/94DanielBrown/roasts-api/internal/database"
)

func UpdateAverages(roastModels database.RoastModels, review database.Review, ratingOperation string) error {
	roast, err := roastModels.GetRoastByPrefix(review.RoastKey)
	if err != nil {
		fmt.Println("error")
		return err
	}
	if roast == nil {
		return errors.New("no roast found")
	}

	var newCount int
	if ratingOperation == "plusCount" {
		newCount = roast.ReviewCount + 1
		roast.OverallRating = ((roast.OverallRating * float64(roast.ReviewCount)) + float64(review.OverallRating)) / float64(newCount)
		roast.MeatRating = ((roast.MeatRating * float64(roast.ReviewCount)) + float64(review.MeatRating)) / float64(newCount)
		roast.PotatoesRating = ((roast.PotatoesRating * float64(roast.ReviewCount)) + float64(review.PotatoesRating)) / float64(newCount)
		roast.VegRating = ((roast.VegRating * float64(roast.ReviewCount)) + float64(review.VegRating)) / float64(newCount)
		roast.GravyRating = ((roast.GravyRating * float64(roast.ReviewCount)) + float64(review.GravyRating)) / float64(newCount)
	} else if ratingOperation == "minusCount" {
		newCount = roast.ReviewCount - 1
		roast.OverallRating = ((roast.OverallRating * float64(roast.ReviewCount)) - float64(review.OverallRating)) / float64(newCount)
		roast.MeatRating = ((roast.MeatRating * float64(roast.ReviewCount)) - float64(review.MeatRating)) / float64(newCount)
		roast.PotatoesRating = ((roast.PotatoesRating * float64(roast.ReviewCount)) - float64(review.PotatoesRating)) / float64(newCount)
		roast.VegRating = ((roast.VegRating * float64(roast.ReviewCount)) - float64(review.VegRating)) / float64(newCount)
		roast.GravyRating = ((roast.GravyRating * float64(roast.ReviewCount)) - float64(review.GravyRating)) / float64(newCount)
	} else {
		return fmt.Errorf("invalid rating operation")
	}

	roast.ReviewCount = newCount

	err = roastModels.UpdateRoast(roast)
	if err != nil {
		return err
	}

	return nil
}
