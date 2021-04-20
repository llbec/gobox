package sixstep

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	//"regexp"
	types "gobox/idcard/govcncode"
)

//SixStep defines a console with six steps
type SixStep struct {
	province *types.Province
	city     *types.City
	county   *types.County
	birthday int
	gender   int
	steps    int
	tree     []*types.Province
}

//NewSixStep return a new sixstep
func NewSixStep(t []*types.Province) *SixStep {
	s := new(SixStep)
	s.tree = t
	s.reset()
	return s
}

func (s *SixStep) reset() {
	s.province = nil
	s.city = nil
	s.county = nil
	s.birthday = 0
	s.gender = 0
	s.steps = 0
	return
}

//Status for join interface console
func (s *SixStep) Status() string {
	return ""
}

//Run start
func (s *SixStep) Run() {
	calc := func() {
		if s.province == nil {
			s.steps = 0
		} else if s.city == nil {
			s.steps = 1
		} else if s.county == nil {
			s.steps = 2
		} else if s.birthday == 0 {
			s.steps = 3
		} else if s.gender == 0 {
			s.steps = 4
		} else {
			//output result
			//fmt.Println(s.province.Name, s.city.Name, s.county.Name, s.birthday, s.gender)
			r, err := s.getresult()
			if err != nil {
				fmt.Printf(err.Error())
			} else {
				title := fmt.Sprintf("%s %s %s %08d ", s.province.Name, s.city.Name, s.county.Name, s.birthday)
				title += func() string {
					if s.gender == 1 {
						return "male"
					}
					return "famle"
				}()
				fmt.Printf(">>>>>>>>>>\n%s\n\tPerson ID is %s\n<<<<<<<<<<\n", title, r)
			}
			s.reset()
		}
	}
	for {
		var input string
		rand.Seed(time.Now().Unix())
		s.showmenu()
		fmt.Scanln(&input)
		if input == "p" {
			s.steps--
		} else if input == "q" {
			fmt.Println("Bye Bye!")
			return
		} else {
			switch s.steps {
			case 0:
				//province
				if input == "" {
					input = strconv.Itoa(rand.Int() % len(s.tree))
				}
				i, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println(err)
					return
				} else if i < 0 && i >= len(s.tree) {
					fmt.Println("Bye Bye!")
					return
				}
				s.province = s.tree[i]
				s.steps++
			case 1:
				//city
				if input == "" {
					input = strconv.Itoa(rand.Int() % len(s.province.Citys))
				}
				i, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println(err)
					return
				} else if i < 0 && i >= len(s.province.Citys) {
					fmt.Println("Bye Bye!")
					return
				}
				s.city = s.province.Citys[i]
				s.steps++
			case 2:
				//county
				if input == "" {
					input = strconv.Itoa(rand.Int() % len(s.city.Countys))
				}
				i, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println(err)
					return
				} else if i < 0 && i >= len(s.city.Countys) {
					fmt.Println("Bye Bye!")
					return
				}
				s.county = s.city.Countys[i]
				s.steps++
			case 3:
				//birthday
				//redate := regexp.MustCompile("^[\\d]{8}$")
				if input == "" {
					input = "19970602"
				}
				i, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println(err)
					return
				}
				s.birthday = i
				s.steps++
			case 4:
				//gender
				if input == "" {
					input = strconv.Itoa((rand.Int() % 2) + 1)
				}
				i, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println(err)
					return
				}
				s.gender = i
				s.steps++
				calc()
			default:
				calc()
			}
		}
	}
}

func (s *SixStep) showmenu() {
START:
	menu := fmt.Sprintf("********** Setp %d **********\nPlease enter your choise:\n", s.steps)
	switch s.steps {
	case 0:
		for i, v := range s.tree {
			menu += fmt.Sprintf("(%d, %s)\t", i, v.Name)
			if (i+1)%5 == 0 {
				menu += "\n"
			}
		}
	case 1:
		if s.province == nil {
			s.steps = 0
			goto START
		}
		for i, v := range s.province.Citys {
			menu += fmt.Sprintf("(%d, %s)\t", i, v.Name)
			if (i+1)%5 == 0 {
				menu += "\n"
			}
		}
	case 2:
		if s.city == nil {
			s.steps = 1
			goto START
		}
		for i, v := range s.city.Countys {
			menu += fmt.Sprintf("(%d, %s)\t", i, v.Name)
			if (i+1)%5 == 0 {
				menu += "\n"
			}
		}
	case 3:
		menu = "********** Setp 3 **********\nPlease enter a birthday(20010101):\n"
	case 4:
		menu += "1 male, 2 famale"
	default:
		if s.province == nil {
			s.steps = 0
			goto START
		} else if s.city == nil {
			s.steps = 1
			goto START
		} else if s.county == nil {
			s.steps = 2
			goto START
		} else if s.birthday == 0 {
			s.steps = 3
			goto START
		} else if s.gender == 0 {
			s.steps = 4
			goto START
		}
		return
	}
	menu += fmt.Sprintf("\n(p, Prev)    (Enter, random)    (others|q, quit)\n")
	fmt.Printf(menu)
}

func (s *SixStep) getresult() (string, error) {
	if s.province == nil {
		return "", fmt.Errorf("Error without province infomation")
	} else if s.city == nil {
		return "", fmt.Errorf("Error without city infomation")
	} else if s.county == nil {
		return "", fmt.Errorf("Error without county infomation")
	} else if s.birthday == 0 {
		return "", fmt.Errorf("Error without birthday infomation")
	} else if s.gender == 0 {
		return "", fmt.Errorf("Error without gender infomation")
	}
	str := fmt.Sprintf("%06d", s.county.GetCode())
	str += fmt.Sprintf("%08d", s.birthday)
	rand.Seed(time.Now().Unix())
	str += fmt.Sprintf("%03d", ((rand.Int()%9)+1)*s.gender)
	r, err := CalcCode(str)
	if err != nil {
		return "", fmt.Errorf("Error length: %s(%d)", str, len(str))
	}
	return str + r, nil
}
