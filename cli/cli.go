package cli

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
)

// LocalCli 本地cli命令
type LocalCli struct {
}

func (cli LocalCli) ExecCommand(command string) {
	logrus.Infof("will send command: %s", command)
	sys := runtime.GOOS
	cmd := &exec.Cmd{}
	if sys == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else if sys == "linux" {
		cmd = exec.Command("sh", "-c", command)
	} else {
		logrus.Warnf("cur sys: %s, not support.", sys)
		return
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("send command: %s error, %+v", command, err)
		return
	}
	logrus.Infof("send command: %s success.", command)
	res := string(output)
	logrus.Infof("res: %s", res)
}
