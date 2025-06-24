package main

import "botvideosaver/internal/logger"

func main() {
	logger.InitLogger()

	logger.Log.Sugar().Info("Starting bot...")

	logger.Log.Sugar().Errorf("This example is not implemented yet. Please refer to the documentation for more information.")

	logger.Log.Sugar().Info("Bot started successfully!")

}
