package main

import "os"
import "bytes"
import "fmt"
import "regexp"
import "strconv"
import "time"
import "math/rand"

type myRegexp struct {
	*regexp.Regexp
}

var dice_expr_regex = myRegexp{regexp.MustCompile(`((?P<rolls>[0-9]*)d)(?P<die>[0-9]+)((?P<op>\+|-)(?P<additive>[0-9]+))?`)}

func main() {
	rand.Seed(time.Now().Unix())

	var explain bool
	for _, cmd := range os.Args[1:] {

		if cmd == "explain" {
			explain = true
			continue
		} else if cmd == "roll" {
			explain = false
			continue
		}

		expr, plan := ComputeExpression(cmd, explain)
		if explain {
			fmt.Printf("%s Rolls {%s} : %d", cmd, plan, expr)
		} else {
			fmt.Printf("%s - %d", cmd, expr)
		}
		fmt.Println("")
	}
}

func ParseInteger(st string, prm map[string]string, def int) int {
	var ret = def
	if value, ok := prm[st]; ok {
		ret, _ = strconv.Atoi(value)
	}
	return ret
}

func ComputeExpression(cmd string, explain bool) (int, string) {
	diceRoll := dice_expr_regex.FindStringSubmatchMap(cmd)
	var buffer bytes.Buffer
	var counter = 0
 	var opAdd = true

	rolls := ParseInteger("rolls", diceRoll, 1)
	die := ParseInteger("die", diceRoll, 0)
	additive := ParseInteger("additive", diceRoll, 0)

	if value, ok := diceRoll["op"]; ok {
		if value == "+" {
			opAdd = true
		} else {
			opAdd = false
		}
	}

	if explain {
		buffer.WriteString("( ")
	}

	die -= 1
	for i := 1; i <= rolls; i++ {
		var newRoll = rand.Intn(die) + 1
		if explain {
			buffer.WriteString(strconv.Itoa(newRoll))
			if i != rolls {
				buffer.WriteString(" + ")
			} else {
				buffer.WriteString(" ")
			}
		}
		counter += newRoll
	}

	if explain {
		buffer.WriteString(")")
	}

	if additive > 0 {
		if opAdd {
			if explain {
				buffer.WriteString(" + ")
				buffer.WriteString(strconv.Itoa(additive))
			}
			counter += additive
		} else {
			if explain {
				buffer.WriteString(" - ")
				buffer.WriteString(strconv.Itoa(additive))
			}
			counter -= additive
		}
	}

	return counter, buffer.String()
}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}
