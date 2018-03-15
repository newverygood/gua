package yi

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	TOTAL  = 49
	UNIT   = 4
	SEQLEN = 6

	YANG = "▅▅▅▅▅"
	YIN  = "▅▅ ▅▅"

	LAOYIN   = 6 // 老阴
	SHAOYIN  = 8 // 少阴
	LAOYANG  = 9 // 老阳
	SHAOYANG = 7 // 少阳
)

type YaoType struct {
	Image string // 爻象
	Text  string // 爻辞
}

type DataType struct {
	Index string          // 卦序
	Name  string          // 卦名
	Text  string          // 卦辞
	Extra string          // 额外信息，如用九、用六
	Short string          // 卦简介
	Desc  string          // 介绍
	Yao   [SEQLEN]YaoType // 六爻
}

var (
	Data map[string]DataType // 易经数据
)

func Load() {
	data, err := Asset("data.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &Data); err != nil {
		log.Fatal(err)
	}
}

////////////////////////////////////////////////////////////
type GuaType struct {
	No [SEQLEN]int // 卜算数字
}

func (g *GuaType) Show() {
	data := Data[g.GetDataIndex()]
	fmt.Printf("%s【卦%s】：%s\n", data.Name, data.Index, data.Text)
	fmt.Println(data.Short)

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"卦象", "爻辞"})
	for _, yao := range data.Yao {
		tb.Append([]string{yao.Image, yao.Text})
	}
	tb.Render()
}

// 变卦
func (g *GuaType) Change() GuaType {
	gc := GuaType{}
	for i, n := range g.No {
		switch n {
		case LAOYIN: // 老阴
			gc.No[i] = LAOYANG
		case SHAOYANG: // 少阳
			gc.No[i] = n
		case SHAOYIN: // 少阴
			gc.No[i] = n
		case LAOYANG: // 老阳
			gc.No[i] = LAOYIN
		}
	}
	return gc
}

func (g *GuaType) Input(args []string) {
	var err error
	for i := 0; i < SEQLEN; i++ { // 自下至上，从0到5
		if g.No[i], err = strconv.Atoi(args[i]); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

// 算卦数古典方法
func (g *GuaType) CalcClassic() {
	rand.Seed(time.Now().UnixNano()) // TODO: 随机数按照日期、卜算内容生成
	for i := 0; i < SEQLEN; i++ {    // 自下至上，从0到5
		g.No[i] = CalcYaoClassic()
	}
	log.Println(g.No)
}

// 算卦数简要方法
func (g *GuaType) CalcSimple(args []string) {
	if len(args) == SEQLEN {
		g.Input(args)
	} else {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < SEQLEN; i++ { // 自下至上，从0到5
			g.No[i] = CalcYaoSimple()
		}
		log.Println(g.No)
	}
}

func (g *GuaType) GetDataIndex() string {
	idx := ""
	for _, n := range g.No {
		idx += strconv.Itoa(n % 2)
	}
	log.Printf("Index [%s]\n", idx)
	return idx
}

// 计算卦象
// 三变之后剩余数字除以4。
func CalcYaoClassic() int {
	b1 := loop(TOTAL)           // 一变
	b2 := loop(TOTAL - b1)      // 二变
	b3 := loop(TOTAL - b1 - b2) // 三变

	m := TOTAL - b1 - b2 - b3
	y := m / UNIT // 爻
	log.Printf("余：%d，爻：%d\n", m, y)

	return y
}

// 变
func loop(total int) int {
	// 天地，随机分组
	l := rand.Intn(total-UNIT*2) + UNIT
	r := total - l

	// 取人，任意一组取一
	if rand.Intn(2) == 0 {
		l -= 1
	} else {
		r -= 1
	}

	// 取天，天数除四取余
	lm := l % UNIT
	if lm == 0 {
		lm = UNIT
	}
	l -= lm

	// 取地，地数除四取余
	rm := r % UNIT
	if rm == 0 {
		rm = UNIT
	}
	r -= rm

	// 取余，取出总数
	m := lm + rm + 1

	log.Printf("总：%d，天：%d，地：%d，取：%d\n", total, l, r, m)
	return m
}

// 计算卦象
// 正面为阳，反面为阴。如果把三枚硬币同时扔下就会有四种组合：
// 二阴一阳、二阳一阴、三阴、三阳。
// 二阴一阳是少阳，二阳一阴是少阴，三阳是老阳，三阴是老阴。
func CalcYaoSimple() int {
	sum := 0
	for i := 0; i < 3; i++ {
		sum += rand.Intn(2)
	}

	switch sum {
	case 0: // 老阴
		return LAOYIN
	case 1: // 少阳
		return SHAOYANG
	case 2: // 少阴
		return SHAOYIN
	case 3: // 老阳
		return LAOYANG
	}
	return 0
}

////////////////////////////////////////////////////////////
// 得位
func (g *GuaType) InPos(n int) bool {
	return true // TODO: 得位算法
}

// 得中
func (g *GuaType) InMid(n int) bool {
	return n == 1 || n == 4
}

func (g *GuaType) Divining() { // TODO: 占卜算法
	var chn []int
	for n, i := range g.No {
		if i == 6 || i == 9 {
			chn = append(chn, n)
		}
	}

	switch len(chn) { // 算出来的六爻当中
	case 0: // 六爻一个都没变，这时用本卦的卦辞来判断吉凶。
	case 1: // 有一个爻是变爻，用本卦变爻的爻辞来判断吉凶。
	case 2: // 有两个爻发生变动，用本卦里这两个变爻的爻辞来判断吉凶，并以位置靠上的那一个爻辞为主。
	case 3: // 有三个变爻，就不能用变爻的爻辞来判断了，得用本卦和变卦的卦辞，以本卦的卦辞为主。
	case 4: // 有四个变爻，这时就用变卦的两个不变爻的爻辞来判断吉凶。
	case 5: // 有五个变爻，用变卦的那一个不变爻的爻辞来判断吉凶。
	case 6: // 有六个变爻，分两种情况。
		if true { // 一是六爻都是阳爻（构成了乾卦），或者六爻都是阴爻（构成了坤卦），那么，
			// 如果是乾卦，就用乾卦“用九”的爻辞判断吉凶，
			// 如果是坤卦，就用坤卦“用六”的爻辞判断吉凶。

		} else { // 二是除了这两种情况之外的其他六爻全变的情况，就用变卦的卦辞来判断吉凶。

		}

	}
	fmt.Println(chn)
}