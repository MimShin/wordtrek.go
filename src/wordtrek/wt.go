package wordtrek

import (
    "fmt"
    "dict"
    "strings"
    "sync"
    "util"
)

type WordTrek struct {
    mu sync.Mutex
    wtns map[string]WTNode
    dict dict.Dict
    rows, cols int
    wordLengths []int
}

func (wt *WordTrek) Solve(tableStr string, wordLengths []int, dict dict.Dict) {

    t := util.StrToTable(tableStr)
    wt.rows, wt.cols = len(t), len(t[0])
    wt.wtns = make(map[string]WTNode) 
    wt.wordLengths = wordLengths;
    wt.dict = dict

	wt.wtns[""] = WTNode{
		words: []string{}, 
		table: t }

    maxLevels := len(wt.wordLengths)
    for i:=0; i<maxLevels; i++ {

        fmt.Printf("Level %d: %d node(s)\n", i, len(wt.wtns))

        workingList := make(map[string]WTNode)
        for key, wtn := range wt.wtns {
            if wtn.level == i {
                workingList[key] = wtn
            }
        }

        for _, wtn := range workingList {
            wt.findWord(wtn)
        }
    }

    for _, wtn := range wt.wtns {
        if wtn.level == maxLevels {
            //wtn.Print();
            fmt.Printf("%s - %d - %c\n", wtn.words, wt.wordLengths, wtn.table)
        }
    }
}

func (wt *WordTrek) findWord(wtn WTNode) {

    var wg sync.WaitGroup

    t := wtn.table
    wg.Add(len(t) * len(t[0]))

    for r:=0; r<len(t); r++ {
        for c:=0; c<len(t[0]); c++ {
            go wt.goFindWordAtRC(&wg, wtn.Clone(), r, c) 
        }
    }
    
    wg.Wait()
}


func (wt *WordTrek) goFindWordAtRC(wg* sync.WaitGroup, wtn WTNode, row int, col int) {
    // fmt.Printf("goFindWordAtRC %c @%d,%d %d\n", wtn.table, row, col);
    defer wg.Done()
    wt.findWordAtRC(wtn, row, col, "")
}


func (wt *WordTrek) findWordAtRC(wtn WTNode, row int, col int, prefix string) {
    // fmt.Printf("findWordsAtRC: %c, %d, %d, %q\n", wtn.table, row, col, prefix)

    t := wtn.table

    ch := t[row][col]

    if ch == '.' {
        return
    } 
    t[row][col] = '.'

    if len(prefix) == wt.wordLengths[wtn.level] - 1 {
        if wt.dict.Look(prefix + string(ch)) {
            wt.addTowtns(WTNode{
                words: append(wtn.words, prefix + string(ch)), 
                table: wtn.table.Clone().DropDown(),
                level: wtn.level + 1 })
        }
        t[row][col] = ch
        return
    }

    for ri := -1; ri<2; ri++ {
        r := row + ri
        if r >= 0 && r < wt.rows {
            for ci := -1; ci<2; ci++ {
                c := col + ci
                if c >= 0 && c < wt.cols {
                    wt.findWordAtRC(wtn, r, c, prefix + string(ch))
                }
            }
        }
    }

    t[row][col] = ch
}

func (wt *WordTrek) addTowtns(wtn WTNode) {

//    fmt.Println("addTowtns", strings.Join(wtn.words, ""))
    wt.mu.Lock()
    defer wt.mu.Unlock()
    wt.wtns[wtn.table.ToStr() + strings.Join(wtn.words, "")] = wtn
}
