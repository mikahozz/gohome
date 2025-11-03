package main

func main() {
	scheduler := CreateSunriseSunsetScheduler()
	scheduler.Start()
	defer scheduler.Stop()

	// Keep the main function running
	select {}
}
