package main

import (
	"fmt"
	"os"
	//"path/filepath"
	"strings"

	"github.com/xntrik/hcltm/pkg/spec"
)

type DfdCommand struct {
	*GlobalCmdOptions
	specCfg       *spec.ThreatmodelSpecConfig
	flagOutDir    string
	flagOutFile   string
	flagOverwrite bool
	flagImageFormat string
}

func (c *DfdCommand) Help() string {
	helpText := `
Usage: hcltm dfd [options] -outdir=<directory> <files>

  Generate Data Flow Diagram PNG/SVG files from existing Threat model HCL files
	(as specified by <files>) 

 -outdir=<directory>
   Directory to output image files. Will create directory if it doesn't exist.
   Either this, or -out, must be set

 -out=<filename>
   Name of output image file. Only the first discovered data_flow_diagram will be converted into an image.
   Either this, or -outdir, must be set, defaults to PNG

Options:

 -config=<file>
   Optional config file

 -overwrite

 -imageformat=<png/svg>
`
	return strings.TrimSpace(helpText)
}

func (c *DfdCommand) Run(args []string) int {

	flagSet := c.GetFlagset("dfd")
	flagSet.StringVar(&c.flagOutDir, "outdir", "", "Directory to output PNG files. Will create directory if it doesn't exist. Either this, or -out, must be set")
	flagSet.StringVar(&c.flagOutFile, "out", "", "Name of output PNG file. Either this, or -outdir, must be set")
	flagSet.BoolVar(&c.flagOverwrite, "overwrite", false, "Overwrite existing files in the outdir. Defaults to false")
	flagSet.StringVar(&c.flagImageFormat, "imageformat", "png", "Set image format for output. Either png or svg. Defaults to PNG")
	
	flagSet.Parse(args)

	if c.flagConfig != "" {
		err := c.specCfg.LoadSpecConfigFile(c.flagConfig)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return 1
		}
	}

	if c.flagOutDir == "" && c.flagOutFile == "" {
		fmt.Printf("You must set an -outdir or -out\n\n")
		fmt.Println(c.Help())
		return 1
	}

	if c.flagOutDir != "" && c.flagOutFile != "" {
		fmt.Printf("You must sent an -outdir or -out, but not both\n\n")
		fmt.Println(c.Help())
		return 1
	}

	if c.flagImageFormat != "png" && c.flagImageFormat != "svg" {
		fmt.Printf("-imageformat flag must be png or svg\n\n")
		fmt.Println(c.Help())
		return 1
	}
	//if c.flagOutFile != "" && filepath.Ext(c.flagOutFile) != ".png" {
	//	fmt.Printf("-out flag must end in .png\n\n")
	//	fmt.Println(c.Help())
	//	return 1
	//}

	if len(flagSet.Args()) == 0 {
		fmt.Printf("Please provide file(s)\n\n")
		fmt.Println(c.Help())
		return 1
	} else {

		// We use outfiles to generate a list of output files to validate whether
		// we're overwriting them or not.
		outfiles := []string{}

		// Find all the .hcl files we're going to parse
		AllFiles := findAllFiles(flagSet.Args())

		// Parse all the identified .hcl files - just to determine output files
		for _, file := range AllFiles {
			tmParser := spec.NewThreatmodelParser(c.specCfg)
			err := tmParser.ParseFile(file, false)
			if err != nil {
				fmt.Printf("Error parsing %s: %s\n", file, err)
				return 1
			}

			for _, tm := range tmParser.GetWrapped().Threatmodels {

				if tm.DataFlowDiagram != nil {
					outfile := outfilePath(c.flagOutDir, tm.Name, file, ".png")

					outfiles = append(outfiles, outfile)
				}

			}
		}

		if c.flagOutFile != "" {
			outfiles = []string{c.flagOutFile}
		}

		// Validating existing files - if we're not overwriting
		if !c.flagOverwrite {
			for _, outfile := range outfiles {
				_, err := os.Stat(outfile)
				if !os.IsNotExist(err) {
					fmt.Printf("'%s' already exists\n", outfile)
					return 1
				}
			}
		}

		if len(outfiles) == 0 {
			fmt.Printf("No Data Flow Diagrams found in provided HCL files\n")
			return 1
		}

		if c.flagOutDir != "" {
			err := createOrValidateFolder(c.flagOutDir, c.flagOverwrite)
			if err != nil {
				fmt.Printf("%s\n", err)
				return 1
			}
		}

		for _, file := range AllFiles {
			tmParser := spec.NewThreatmodelParser(c.specCfg)
			err := tmParser.ParseFile(file, false)
			if err != nil {
				fmt.Printf("Error parsing %s: %s\n", file, err)
				return 1
			}

			for _, tm := range tmParser.GetWrapped().Threatmodels {
				if tm.DataFlowDiagram != nil {
					if c.flagOutFile != "" {
						if c.flagImageFormat == "png" {
							err = tm.GenerateDfdPng(c.flagOutFile)	
						} else if c.flagImageFormat == "svg" {
							err = tm.GenerateDfdSvg(c.flagOutFile)
						}
						if err != nil {
							fmt.Printf("Error generating DFD: %s\n", err)
							return 1
						}

						fmt.Printf("Successfully created '%s'\n", c.flagOutFile)
						break
					} else {
						if c.flagImageFormat == "png" {
							err = tm.GenerateDfdPng(outfilePath(c.flagOutDir, tm.Name, file, ".png"))
						} else if c.flagImageFormat == "svg" {
							err = tm.GenerateDfdSvg(outfilePath(c.flagOutDir, tm.Name, file, ".svg"))
						}
						if err != nil {
							fmt.Printf("Error generating DFD: %s\n", err)
							return 1
						}

						if c.flagImageFormat == "png" {
							fmt.Printf("Successfully created '%s'\n", outfilePath(c.flagOutDir, tm.Name, file, ".png"))
						} else if c.flagImageFormat == "svg" {
							fmt.Printf("Successfully created '%s'\n", outfilePath(c.flagOutDir, tm.Name, file, ".svg"))
						}

					}
				}
			}
		}

	}

	return 0
}

func (c *DfdCommand) Synopsis() string {
	return "Generate Data Flow Diagram image files from existing HCL threatmodel file(s)"
}
