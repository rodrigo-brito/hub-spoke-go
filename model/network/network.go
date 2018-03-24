package network

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rodrigo-brito/hub-spoke-go/util/log"
)

// InputData store all network data
type NetworkData struct {
	Size             int
	ScaleFactor      float64
	InstallationCost []float64
	Distance         [][]float64
	Flow             [][]float64
}

// parseLine validate a line of values and parse to a float array
func parseLine(line string) (bool, []float64, error) {
	values := strings.Split(strings.TrimSpace(line), " ")
	numbers := make([]float64, 0)

	for _, value := range values {
		if len(value) == 0 {
			continue
		}

		number, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			log.Error(err)
			return false, nil, err
		}

		numbers = append(numbers, number)
	}

	if len(numbers) > 0 {
		return true, numbers, nil
	}
	return false, numbers, nil
}

// nextLine validate and return the next valid line
func nextLine(scanner *bufio.Scanner) (bool, []float64, error) {
	if ok := scanner.Scan(); !ok {
		return false, nil, fmt.Errorf("unexpected end of file")
	}

	line := scanner.Text()
	if ok, values, err := parseLine(line); ok {
		return true, values, nil
	} else if err != nil {
		return false, nil, err
	}

	// return the next valid line
	return nextLine(scanner)
}

// FromFile read a input file and generate the network data
func FromFile(fileName string) (*NetworkData, error) {
	data := new(NetworkData)

	file, err := os.Open(fileName)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)

	// First line: network size
	if ok, line, err := nextLine(sc); ok {
		data.Size = int(line[0])
	} else if err != nil {
		return nil, err
	}

	// Second line: scale factor
	if ok, line, err := nextLine(sc); ok {
		data.ScaleFactor = line[0]
	} else if err != nil {
		return nil, err
	}

	// Hub installation cost
	data.InstallationCost = make([]float64, data.Size)
	for i := 0; i < data.Size; i++ {
		if ok, line, err := nextLine(sc); ok {
			data.InstallationCost[i] = line[0]
		} else if err != nil {
			return nil, err
		}
	}

	// Distance between nodes
	data.Distance = make([][]float64, data.Size, data.Size)
	for i := 0; i < data.Size; i++ {
		if ok, line, err := nextLine(sc); ok {
			if len(line) != data.Size {
				return nil, fmt.Errorf("distance matrix should have size %d", data.Size)
			}
			data.Distance[i] = line
		} else if err != nil {
			return nil, err
		}
	}

	// Flow between nodes
	data.Flow = make([][]float64, data.Size, data.Size)
	for i := 0; i < data.Size; i++ {
		if ok, line, err := nextLine(sc); ok {
			if len(line) != data.Size {
				return nil, fmt.Errorf("flow matrix should have dimension %dx%d", data.Size, data.Size)
			}
			data.Flow[i] = line
		} else if err != nil {
			return nil, err
		}
	}

	return data, nil
}