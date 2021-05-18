package gocv

import (
	"fmt"

	"gocv.io/x/gocv"
)

// VidsrcGOCV is a concrete implementation of the Vidsrc abstract interface
type VidsrcGOCV struct {
	Camstr string
}

// StreamVideo takes a device string, starts the video stream and
// starts pumping single JPEG images from the camera stream.
//
// TODO: Change this to Camera and create an interface that
// is sufficient for video files and imagnes.
func (cam *Camera) PumpVideo() (frames <-chan *gocv.Mat) {
	var err error

	// Do not try to restart the video when it is already running.
	if cam.Recording {
		log.Error().Msg("camera already recording")
		return nil
	}

	// Create the channel we are going to pump frames through
	// frameQ := make(chan *gocv.Mat)
	frameQ := make(chan *redeye.Frame)

	// go function opens the webcam and starts reading from device, coping frames
	// to the frameQ processing channel
	go func() {

		defer log.Println("cameraid: ", cam.Camstr, " recording: ", cam.Recording, " Stop StreamVideo")

		// Open the camera (capture device)
		var cap *gocv.VideoCapture
		camstr := GetCamstr(cam.Camstr)
		log.Println("Opening Video with camstr: ", camstr, "Opening VideoCapture")

		cap, err = gocv.OpenVideoCapture(camstr)
		if err != nil {
			log.Fatal().Msg("failed to open video capture device")
			return
		}
		defer cap.Close()

		log.Info().Msg("Camera streaming  ...")

		// Only a single static image will be in the system at a given time.
		img := gocv.NewMat()

		// as long as cam.recording is true we will capture images and send
		// them into the image pipeline. We may recieve a REST or MQTT request
		// to stop recording, in that case the cam.recording will be set to
		// false and the recording will stop.
		cam.Recording = true
		for cam.Recording {

			// read a single raw image from the cam.
			if ok := cap.Read(&img); !ok {
				log.Info().Msg("device closed, turn recording off")
				cam.Recording = false
			}

			// if the image is empty, there will be no sense continueing
			if img.Empty() {
				continue
			}

			// sent the frame to the frame pipeline (should we send by )
			fmt.Printf("frame %+v\n", img)
			frameQ <- &img
		}
		log.Info().Bool("recording", cam.Recording).Msg("Video loop exiting ...")
	}()

	// return the frame channel, our caller will pass it to the reader
	return frameQ
}
