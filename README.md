# Create Test Videos

Create Test videos at given resolutions and colors

## Usage

This `go` package is essentially a shell wrapper. You need to ensure both `imagemagick` and `ffmpeg` are installed and in your $PATH before running this file

### Create YAML Files

The colors and resolutions are read from YAML files by default. Create a `colors.yml` and `resolutions.yml`. You can copy the included sample files as a good starting point, and edit them as necessary.

### Get the binary

#### Option 1: Download the binary

Download the latest built version from [releases](https://github.com/ehowe/create_test_videos/releases)

#### Option 2: Build the Source

```
git clone git@github.com:ehowe/create_test_videos
cd create_test_videos
make
```

### Run it

```
Usage of ./create-test-videos:
  -c, --colors string        Required: Path to colors YAML file
  -d, --dry-run              Dont actually create anything
  -h, --help                 Show help text
  -o, --output-dir string    Output path for generated videos (default `pwd`)
  -r, --resolutions string   Required: Path to resolutions YAML file
  -v, --verbose              Verbose output
```

In practice, this looks like this:

```
./create-test-videos -c colors.yml -r resolutions.yml -o <some directory> -v
```

If you use the sample files that I provided, this will take a _long_ time. The larger the resolutions, the longer it takes. Good luck.
