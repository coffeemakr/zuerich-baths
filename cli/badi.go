package main

import (
	"fmt"
	"os"
	"strings"

	baths "github.com/coffeemakr/zuerich-baths"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "[Name]",
	Short: "Listed Frei-Badis auf.",
	Args:  cobra.MaximumNArgs(1),
	Run:   printBaths,
}

var printBathDetailsCommand = &cobra.Command{
	Use:   "zeig <Name>",
	Short: "Zeigt details zu einer Badi.",
	Args:  cobra.ExactArgs(1),
	Run:   printBathDetails,
}

func init() {
	rootCommand.AddCommand(printBathDetailsCommand)
}

func filter(vs []*baths.ZuerichBaths, f func(*baths.ZuerichBaths) bool) []*baths.ZuerichBaths {
	vsf := make([]*baths.ZuerichBaths, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func filterBathsByName(elements []*baths.ZuerichBaths, name string) []*baths.ZuerichBaths {
	name = strings.ToLower(name)
	return filter(elements, func(bath *baths.ZuerichBaths) bool {
		return strings.Contains(strings.ToLower(bath.Name), name)
	})
}

func printBaths(cmd *cobra.Command, args []string) {
	baths, err := baths.GetBaths()
	if err != nil {
		fmt.Printf("Fehler: Konnte nicht geladen werden.\nMehr Infos: %s", err)
		os.Exit(1)
	}

	if len(args) != 0 {
		baths = filterBathsByName(baths, args[0])
	}
	for _, bath := range baths {
		fmt.Printf("%-30s: %2.1f°C\n", bath.Name, bath.WaterTemperature)
	}
}

func printBathDetails(cmd *cobra.Command, args []string) {
	baths, err := baths.GetBaths()
	if err != nil {
		fmt.Printf("Fehler: Konnte nicht geladen werden.\nMehr Infos: %s", err)
		os.Exit(1)
	}

	baths = filterBathsByName(baths, args[0])
	slalom := "~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°-.~¨°+"
	for _, bath := range baths {
		fmt.Printf("\033[34m(%s\n ) \033[1;39m%-76s\033[22;34m |\n(%s\033[0m\n\033[34m )\033[0m\n", slalom, bath.Name, slalom)
		fmt.Printf("\033[34m( \033[0m Wasser Temperatur: %.0f°C\n", bath.WaterTemperature)
		fmt.Printf("\033[34m )\033[0m Öffnungszeiten:    %s\n", bath.OpenClosedTextPlain)
		fmt.Printf("\033[34m( \033[0m Webseite:          %s\n", bath.URLPage)
		fmt.Printf("\033[34m )%s\033[0m\n", slalom)
	}
}

func main() {
	rootCommand.Execute()
}
