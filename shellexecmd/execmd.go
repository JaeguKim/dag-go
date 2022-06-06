package shellexecmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// TODO exec.Command 에서 파일 위치에 대한 검사를 하지만 오류를 리턴하지는 않는다.
func ScriptRunner(s string) (*exec.Cmd, io.Reader) {
	cmd := exec.Command(s)

	// StdoutPipe 쓰면 Run 및 기타 Run 을 포함한 method 를 쓰면 에러난다.
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicf("Error stdout pipe for Cmd: %v", err)
	}

	return cmd, r
}

// TODO
// https://yourbasic.org/golang/multiline-string/
// string으로 가지고 와서 shell command 만 오는 것이 아니라, \t\n 같은 escape letter 들도 온다. 물론 실행에는 문제 없지만
// 잠재적인 오류를 없애기 위해서 shell command 만 가지고 와야 하지 않을까??
func ScriptRunnerString(s string) (*exec.Cmd, io.Reader) {
	cmd := exec.Command("/bin/sh", "-c", s)

	// StdoutPipe 쓰면 Run 및 기타 Run 을 포함한 method 를 쓰면 에러난다.
	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Panicf("Error stdout pipe for Cmd: %v", err)
	}

	return cmd, r
}

func StartThenWait(cmd *exec.Cmd) {
	go func(cmd *exec.Cmd) {
		if cmd != nil {
			if err := cmd.Start(); err != nil {
				log.Printf("Error starting Cmd: %v", err)
				return
			}
			if err := cmd.Wait(); err != nil {
				log.Printf("Error waiting for Cmd: %v", err)
				return
			}
		}
	}(cmd)
}

func Reply(i io.Reader) <-chan string {
	r := make(chan string, 1)

	go func() {
		defer close(r)
		scan := bufio.NewScanner(i)

		for {
			b := scan.Scan()
			if b != true {
				if scan.Err() == nil {
					// grpc 에서는 스트림을 닫아버리자.
					r <- "FINISHED"
					break
				}
				log.Println(scan.Err())
				r <- "ERRORS"
				break
			}

			s := scan.Text()
			r <- s
		}
	}()
	return r
}

func PrintOutput(ch <-chan string) {
	for m := range ch {
		if strings.Contains(m, "FINISHED") {
			log.Println("Exit Ok")
			return
		}
		if strings.Contains(m, "ERRORS") {
			log.Println("Exit Error")
			return
		}
		fmt.Println(">", m)
	}
}

func Runner(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return false
	}
	cmd, r := ScriptRunnerString(s)
	StartThenWait(cmd)
	ch := Reply(r)

	PrintOutput(ch)

	return true
}
