package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/CoatHanger139/Vremenska-Postaja/vreme"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "weatherStation",
		Usage: "Weather station displays Temperature",
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:  "minTemp",
				Usage: "najnižja dovoljena temperatura meritve (privzeto -50 °C)",
				Value: -50.0,
			},
			&cli.Float64Flag{
				Name:  "maxTemp",
				Usage: "najvišja dovoljena temperatura meritve (privzeto 50 °C)",
				Value: 50.0,
			},
			&cli.IntFlag{
				Name: "stMeritev",
				Usage: " najmanjše število meritev, potrebnih za izračun statistike (privzeto 3)",
				Value: 3,
			},
		},
		Action: func( ctx context.Context, cmd *cli.Command ) error {
			minTemp := cmd.Float64( "minTemp" )
			maxTemp := cmd.Float64( "maxTemp" )
			stMeritev := cmd.Int( "stMeritev" )
			return run( minTemp, maxTemp, stMeritev )
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Print( err.Error() )
	}
}

func run( minTemp, maxTemp float64, stMeritev int ) error {
	stationMap := map[ string ]vreme.Postaja{
		"KP01": vreme.UstvariPostajo( "Koper Luka"        , "Koper"    , minTemp, maxTemp ),
		"KR01": vreme.UstvariPostajo( "Kredarica" 		  , "Triglav"  , minTemp, maxTemp ),
		"LJ01": vreme.UstvariPostajo( "Ljubljana Bežigrad", "Ljubljana", minTemp, maxTemp ),
		"MB01": vreme.UstvariPostajo( "Maribor Tabor"     , "Maribor"  , minTemp, maxTemp ),
	}

	var quit chan struct{}
	quit = make( chan struct{}, 0 )
	go stop( quit )
	i := 0
	for {
		i++
		for key := range stationMap {
			temp := rand.Float64() * 150 -100
			vreme.DodajMeritev( stationMap, key, temp )
			// vreme.IzpisMeritev( stationMap )
		}
		if i % stMeritev == 0 {
			vreme.IzpisPorocila( stationMap )
		}

		select {
			case <-quit:
				return nil
			default:
				time.Sleep( 250 * time.Millisecond )
				continue
		}
	}
}

func stop( ch chan struct{} ) {
	fmt.Scanln()
	close( ch )
}