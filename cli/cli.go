package cli

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
	"strings"
)

// LocalCli 本地cli命令
type LocalCli struct {
}

func (cli LocalCli) ExecCommand(command string, arg ...string) (string, error) {
	sys := runtime.GOOS
	fullCommand := command + " " + strings.Join(arg, " ")
	logrus.Infof("will send command: %s in %s", fullCommand, sys)
	cmd := &exec.Cmd{}
	if sys == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else if sys == "linux" || sys == "darwin" {
		cmd = exec.Command(command, arg...)
	} else {
		logrus.Warnf("cur sys: %s, not support.", sys)
		return "", errors.New(fmt.Sprintf("%s not support", sys))
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("send command: %s error, %+v", fullCommand, err)
		return "", err
	}
	logrus.Infof("send command: %s success.", fullCommand)
	res := string(output)
	logrus.Infof("res: %s", res)
	return res, nil
}
