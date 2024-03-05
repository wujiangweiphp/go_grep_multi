package scan_project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type SmartScan struct {
	RootPath       string
	IncludePathStr string
	ExcludePathStr string
	NeedLineCount  bool
	NeedUseRegexp  bool
	OnlyFile       bool
	MatchContent   []string
	IgnoreCate     bool
}

type MatchRes struct {
	MaStr string
	Line  int
}

type MR []MatchRes

func (l MR) Len() int {
	return len(l)
}
func (l MR) Less(i, j int) bool {
	return l[i].Line < l[j].Line
}

func (l MR) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (s *SmartScan) Scan() error {
	totalLine := 0
	err := filepath.Walk(s.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if s.IncludePathStr != "" && !strings.Contains(path, s.IncludePathStr) {
			return nil
		}

		if s.ExcludePathStr != "" && strings.Contains(path, s.ExcludePathStr) {
			return nil
		}

		if s.OnlyFile {
			totalLine, err = s.readFile(path, totalLine)
		} else {
			if s.MatchContent[0] == "use" {
				s.noUseScan(path)
				return nil
			}
			tempLine := 0
			flagCount := 0
			printFlag := false
			var resStr MR
			for colorKey, matchStr := range s.MatchContent {
				flag := false
				var temp MR
				tempLine, flag, temp, _ = s.readByLine(path, matchStr, 0, colorKey)
				// fmt.Println(totalLine, flag, temp)
				if flag {
					flagCount++
					resStr = append(resStr, temp...)
				}
			}
			totalLine += tempLine
			l := len(s.MatchContent)
			if l > 1 && flagCount == l {
				printFlag = true
			} else if l == 1 {
				printFlag = true
			}
			if len(resStr) > 0 && printFlag {
				fmt.Println("file:", path)
				sort.Sort(resStr)
				for _, v := range resStr {
					fmt.Println(strconv.Itoa(v.Line) + ": " + v.MaStr)
				}
				fmt.Println("---------------------------")
			}
		}

		return err
	})
	if s.NeedLineCount {
		fmt.Println("扫描总行数：", totalLine)
	}
	return err
}

func (s *SmartScan) noUseScan(filename string) {
	// 1. 先扫描出use
	ms := s.onlyFileMatch(filename, "use", true)
	if len(ms) == 0 {
		return
	}
	var mk []struct {
		MatchStr string
		Line     int
		IsExists bool
		SubStr   string
	}
	for _, v := range ms {
		tmp := strings.ReplaceAll(v.MaStr, "use", "")
		tmp = strings.ReplaceAll(tmp, " ", "")
		tmp = strings.ReplaceAll(tmp, ";", "")
		tmpArr := strings.Split(tmp, "\\")
		l := len(tmpArr)
		mk = append(mk, struct {
			MatchStr string
			Line     int
			IsExists bool
			SubStr   string
		}{MatchStr: v.MaStr, Line: v.Line, IsExists: false, SubStr: tmpArr[l-1]})
		// fmt.Println(tmpArr[l-1])
	}
	// 2. 再根据解析 扫描使用处
	if len(mk) == 0 {
		return
	}
	printFlag := false
	for k, vv := range mk {
		mss := s.onlyFileMatch(filename, vv.SubStr, false)
		if len(mss) == 1 {
			mk[k].IsExists = true
			printFlag = true
		}
	}

	if len(mk) > 0 && printFlag {
		fmt.Println("file:", filename)
		for _, v := range mk {
			if v.IsExists {
				fmt.Println(strconv.Itoa(v.Line) + ": " + v.MatchStr)
			}
		}
		fmt.Println("---------------------------")
	}
}

func (s *SmartScan) onlyFileMatch(filename, matchStr string, isPrefix bool) []MatchRes {
	f, e := os.Open(filename)
	if e != nil {
		return nil
	}
	defer f.Close()
	var ms []MatchRes
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		lineText := scanner.Text()
		ok := false
		if isPrefix {
			ok = strings.HasPrefix(lineText, matchStr)
		} else {
			ok = strings.Contains(lineText, matchStr)
		}
		if ok {
			ms = append(ms, MatchRes{lineText, line})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	return ms
}

func (s *SmartScan) readByLine(filename, matchStr string, totalLine, colorKey int) (int, bool, []MatchRes, error) {
	f, e := os.Open(filename)
	if e != nil {
		return totalLine, false, nil, e
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	var resStr []MatchRes
	for scanner.Scan() {
		totalLine++
		line++
		lineText := scanner.Text()
		if ok, res := s.matchText(lineText, matchStr, colorKey); ok {
			resStr = append(resStr, MatchRes{res, line})
		}
	}
	if err := scanner.Err(); err != nil {
		return totalLine, false, nil, err
	}
	flag := false
	if len(resStr) > 0 {
		flag = true
	}
	return totalLine, flag, resStr, nil
}

func (s *SmartScan) matchText(lineText, matchStr string, colorKey int) (bool, string) {
	var reg *regexp.Regexp
	if s.NeedUseRegexp {
		reg = regexp.MustCompile(matchStr)
	}
	colorMap := map[int]string{
		0: "\u001B[31m%s\u001B[0m",
		1: "\u001B[32m%s\u001B[0m",
		2: "\u001B[33m%s\u001B[0m",
		3: "\u001B[34m%s\u001B[0m",
		4: "\u001B[35m%s\u001B[0m",
		5: "\u001B[36m%s\u001B[0m",
	}
	flag := false
	if s.NeedUseRegexp {
		repStr := fmt.Sprintf(colorMap[colorKey], "$0")
		r := reg.FindString(lineText)
		flag = len(r) > 0
		lineText = reg.ReplaceAllString(lineText, repStr)
	} else {
		repStr := fmt.Sprintf(colorMap[colorKey], matchStr)
		if s.IgnoreCate {
			repStr = fmt.Sprintf(colorMap[colorKey], "$0")
		}
		flag = s.ContainsIgnoreCase(lineText, matchStr)
		lineText = s.ReplaceIgnoreCase(lineText, matchStr, repStr)
	}
	return flag, lineText
}

func (s *SmartScan) readFile(lineText string, totalLine int) (int, error) {
	var reg *regexp.Regexp
	matchStr := s.MatchContent[0]
	if s.NeedUseRegexp {
		reg = regexp.MustCompile(matchStr)
	}
	flag := false
	if s.NeedUseRegexp {
		r := reg.FindString(lineText)
		flag = len(r) > 0
		lineText = reg.ReplaceAllString(lineText, "\033[31m$0"+"\033[0m")
	} else {
		flag = s.ContainsIgnoreCase(lineText, matchStr)
		lineText = strings.ReplaceAll(lineText, matchStr, "\033[31m"+matchStr+"\033[0m")
	}
	if flag {
		totalLine++
		fmt.Println(lineText)
		fmt.Println("---------------------------")
	}
	return totalLine, nil
}

func (s *SmartScan) ContainsIgnoreCase(str, substr string) bool {
	if s.IgnoreCate {
		return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
	}
	return strings.Contains(str, substr)
}

func (s *SmartScan) ReplaceIgnoreCase(str, oldStr, newStr string) string {
	if s.IgnoreCate {
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(oldStr))
		return re.ReplaceAllString(str, newStr)
	}
	return strings.ReplaceAll(str, oldStr, newStr)
}
