package utilCmd

import (
	"github.com/hilaoyu/go-utils/utils"
	"os/exec"
	"strings"
)

func buildExec(name string, args ...string) *exec.Cmd {

	cmdExec := exec.Command(name, args...)
	if "windows" == utils.RunningOs("windows") {
		cmdExec = exec.Command("cmd", "/c", cmdExec.String())
	}

	return cmdExec

}

func RunCommand(wait bool, name string, args ...string) (string, error) {

	cmdExec := buildExec(name, args...)
	///fmt.Println("Running cmd:" + cmdExec.String())
	if wait {
		result, err := cmdExec.CombinedOutput()
		/*if err != nil {
			return result, err
		}*/
		return strings.TrimSpace(string(result)), err
	} else {
		err := cmdExec.Start()
		return "", err
	}
}

// 根据进程名判断进程是否运行
func HasRunning(serverName string, threads ...int) bool {

	checkThreads := 1
	if len(threads) > 0 {
		checkThreads = threads[0]
	}
	if checkThreads <= 0 {
		checkThreads = 1
	}
	cmd := "sh"
	cmdArgs := []string{"-c", "ps -ef  --cols 1000 | grep \"" + serverName + "\" | grep -v grep"}
	//filepath.Join(utils.GetSelfPath(),"hasRunSuricata.sh")
	if "windows" == utils.RunningOs("windows") {
		cmd = `tasklist /fi "imagename eq ` + serverName + `" /FO TABLE /NH`
		cmdArgs = []string{}
	}

	result, err := RunCommand(true, cmd, cmdArgs...)

	if err != nil {
		return false
	}
	if "windows" == utils.RunningOs("windows") {
		return strings.HasPrefix(result, "信息: 没有运行的任务匹配指定标准")
	}

	results := strings.Split(result, "\n")
	return len(results) >= checkThreads
}
