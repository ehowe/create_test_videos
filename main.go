package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/samber/lo"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Color struct {
	Name string `yaml:"name"`
	Hex  string `yaml:"hex"`
}

type Resolution struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

var colors []Color
var resolutions []Resolution

func main() {
	path, err := os.Getwd()

	if err != nil {
		log.Fatal("Cannot get current working director")
	}

	var colorsFlag *string = flag.StringP("colors", "c", "", "Required: Path to colors YAML file")
	var resolutionsFlag *string = flag.StringP("resolutions", "r", "", "Required: Path to resolutions YAML file")
	var outputFlag *string = flag.StringP("output-dir", "o", path, "Output path for generated videos")
	var dryRunFlag *bool = flag.BoolP("dry-run", "d", false, "Dont actually create anything")
	var verboseFlag *bool = flag.BoolP("verbose", "v", false, "Verbose output")
	var helpFlag *bool = flag.BoolP("help", "h", false, "Show help text")

	Usage := func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	DryRunOrExecute := func(command string) {
		if *dryRunFlag {
			fmt.Printf("DRY RUN: %s\n", command)
		} else {
			if *verboseFlag {
				fmt.Println(command)
			}

			cmd := exec.Command("bash", "-c", command)

			if err := cmd.Start(); err != nil {
				log.Printf("Failed to start cmd: %v", err)
				return
			}

			if err := cmd.Wait(); err != nil {
				log.Printf("Cmd returned error: %v", err)
			}
		}
	}

	flag.Parse()

	if *helpFlag {
		Usage()

		os.Exit(0)
	}

	if *colorsFlag == "" || *resolutionsFlag == "" {
		Usage()

		os.Exit(0)
	}

	colorsContents, err := os.ReadFile(*colorsFlag)

	if err != nil {
		log.Fatalf("Error reading colors file %s", *colorsFlag)
	}

	resolutionsContents, err := os.ReadFile(*resolutionsFlag)

	if err != nil {
		log.Fatalf("Error reading resolutions file %s", *resolutionsFlag)
	}

	err = yaml.Unmarshal(colorsContents, &colors)

	if err != nil {
		log.Fatalf("Error unmarshaling colors yaml")
	}

	err = yaml.Unmarshal(resolutionsContents, &resolutions)

	if err != nil {
		log.Fatalf("Error unmarshaling resolutions yaml")
	}

	lo.ForEach(resolutions, func(resolution Resolution, _ int) {
		resolutionName := fmt.Sprintf("%dx%d", resolution.Width, resolution.Height)
		fullPath := filepath.Join(*outputFlag, resolutionName)

		lo.ForEach(colors, func(color Color, index int) {
			imageName := fmt.Sprintf("%s-%s.jpg", resolutionName, color.Name)
			videoName := fmt.Sprintf("%s-%s.mov", resolutionName, color.Name)
			imagePath := filepath.Join(fullPath, imageName)
			videoPath := filepath.Join(fullPath, videoName)

			if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("creating %s image with resolution of %s\n", color.Name, resolutionName)

				DryRunOrExecute(fmt.Sprintf("convert -size %s -gravity center -background '#%s' -fill white -font Arial -pointsize 72 label:%s %s", resolutionName, color.Hex, resolutionName, imagePath))
			} else {
				if *verboseFlag {
					fmt.Printf("Skipping %s because it already exists\n", imagePath)
				}
			}

			if _, err := os.Stat(videoPath); errors.Is(err, os.ErrNotExist) {
				fmt.Printf("creating %s video with resolution of %s from %s\n", color.Name, resolutionName, imagePath)

				DryRunOrExecute(fmt.Sprintf("ffmpeg -r 1/10 -i %s -c:v libx264 -pix_fmt yuv420p %s", imagePath, videoPath))
			} else {
				if *verboseFlag {
					fmt.Printf("Skipping %s because it already exists\n", videoPath)
				}
			}

			var nextColor Color

			if index == len(colors)-1 {
				nextColor = colors[0]
			} else {
				nextColor = colors[index+1]
			}

			fmt.Println(index, len(colors), nextColor)

			transitionName := fmt.Sprintf("%s-%s-to-%s.mov", resolutionName, color.Name, nextColor.Name)
			transitionImagePath := filepath.Join(fullPath, fmt.Sprintf("%s-%s.jpg", resolutionName, nextColor.Name))
			transitionPath := filepath.Join(fullPath, transitionName)

			if _, err := os.Stat(transitionPath); errors.Is(err, os.ErrNotExist) {
				if *verboseFlag {
					fmt.Printf("creating %s transition video with resolution of %s\n", transitionName, resolutionName)
				}

				DryRunOrExecute(fmt.Sprintf("ffmpeg -loop 1 -t 10 -i %s -loop 1 -t 10 -i %s  -filter_complex '[0][1]xfade=transition=diagbr:duration=10' -c:v libx264 -pix_fmt yuv420p %s", imagePath, transitionImagePath, transitionPath))
			} else {
				if *verboseFlag {
					fmt.Printf("Skipping %s because it already exists\n", transitionPath)
				}
			}
		})
	})
}
