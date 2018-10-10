package main

import (
    "fmt"
    "bufio"
    "os"
    "log"
    "strconv"
    "time"
    "sync"
)

type coord struct {
    x int
    y int
}

func print_grid(grid [][]uint8) {
    for i, row := range grid {
        for j, _ := range row {
            fmt.Printf("%d ", grid[j][i])
        }
        fmt.Printf("\n")
    }
    fmt.Printf("\n")
}

func get_num_neighbors(grid [][]uint8, p coord) int {
    n := len(grid)
    num_neighbors := 0

    //abs_neighbor_positions := make([]coord, 8)
    rel_neighbor_positions := []coord{
        coord{-1, -1},
        coord{ 0, -1},
        coord{ 1, -1},
        coord{-1,  0},
        coord{ 1,  0},
        coord{-1,  1},
        coord{ 0,  1},
        coord{ 1,  1}}

    for _, rnp := range rel_neighbor_positions {
        anp := add_coords(rnp, p)
        anp.x = _mod(anp.x, n)
        anp.y = _mod(anp.y, n)
        //fmt.Printf("%d, %d\n", anp.x, anp.y)
        if grid[anp.x][anp.y] == 1 {
            num_neighbors++
        }
    }
    return num_neighbors
}

func _mod(a int, b int) int {
    var r int

    if a%b >= 0 {
        r = a%b
    } else {
        r = a%b + b
    }
    return r
}

func add_coords(a coord, b coord) coord {
    r := coord{a.x+b.x, a.y+b.y}
    return r
}

func update(grid [][]uint8) [][]uint8 {
    n := len(grid)
    next_grid := make([][]uint8, n)
    for i:= range grid {
        next_grid[i] = make([]uint8, n)
    }

    var wg sync.WaitGroup

    for i, _ := range grid {
        wg.Add(1)
        go _update_row(grid, next_grid, i, &wg)
    }

    wg.Wait()
    return next_grid
}

func _update_row(grid [][]uint8, next_grid [][]uint8, i int, wg *sync.WaitGroup) {
    row := grid[i]

    for j,_ := range row {
        num_neighbors := get_num_neighbors(grid, coord{i, j})
        switch grid[i][j] {
        case 1:
            switch {
            case num_neighbors < 2:
                next_grid[i][j] = 0
            case num_neighbors == 2 || num_neighbors == 3:
                next_grid[i][j] = 1
            case num_neighbors > 3:
                next_grid[i][j] = 0
            }
        case 0:
            if num_neighbors == 3 {
                next_grid[i][j] = 1
            }
        }
    }
    defer wg.Done()
}

func main() {

    const n int = 10

    grid := make([][]uint8, n)
    for i:= range grid {
        grid[i] = make([]uint8, n)
    }

    file, err := os.Open("glider")
    if err != nil {
        log.Fatal(err)
    }

    p := coord{0,0}

    scanner := bufio.NewScanner(file)

    j := 0
    for scanner.Scan() {
        for i, c := range scanner.Text() {
            x, _ := strconv.Atoi(string(c))
            grid[p.x+i][p.y+j] = uint8(x)
        }
        j++
    }
    file.Close()

    start := time.Now()
    for i:=0;i<1000;i++{
        grid = update(grid)
    }
    elapsed := time.Since(start)
    fmt.Printf("%d \n", elapsed/1000)

}
