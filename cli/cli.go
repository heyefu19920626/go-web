package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

var WindowsOverSign = []string{">"}
var LinuxOverSign = []string{"#"}

// LocalCli 本地cli命令
type LocalCli struct {
	//结束符
	overSign []string
	//系统类型
	sys string
	//执行命令的通道
	cmd *exec.Cmd
	//命令输入
	in io.WriteCloser
	// 读取标准输出
	standReader *bufio.Reader
	//读取错误输出
	errReader *bufio.Reader
	// 标准结果
	standRes strings.Builder
	// errorRes
	errorRes strings.Builder
}

func NewLocalCli(overSign []string) (*LocalCli, error) {
	localCli := LocalCli{}
	localCli.overSign = overSign
	localCli.sys = getSys()
	if localCli.sys == "windows" {
		localCli.cmd = exec.Command("cmd")
	} else if localCli.sys == "linux" || localCli.sys == "darwin" {
		localCli.cmd = exec.Command("/bin/bash")
	} else {
		logrus.Warnf("cur sys: %s, not support.", localCli.sys)
		return nil, errors.New(fmt.Sprintf("%s not support", localCli.sys))
	}
	logrus.Infof("create cmd finish")
	in, _ := localCli.cmd.StdinPipe()
	outStand, _ := localCli.cmd.StdoutPipe()
	outErr, _ := localCli.cmd.StderrPipe()
	logrus.Infof("start cmd")
	if err := localCli.cmd.Start(); err != nil {
		logrus.Errorf("create cmd error! %+v", err)
		return nil, errors.New("create cmd error")
	}
	logrus.Infof("start cmd finish")
	localCli.standReader = bufio.NewReader(outStand)
	localCli.errReader = bufio.NewReader(outErr)
	localCli.standRes = strings.Builder{}
	localCli.errorRes = strings.Builder{}
	localCli.in = in
	go localCli.readOutput(localCli.standReader, &localCli.standRes)
	go localCli.readOutput(localCli.errReader, &localCli.errorRes)
	logrus.Infof("start read login info")
	start, _ := localCli.GetResult()
	logrus.Infof("get login info: %s", start)
	return &localCli, nil
}

func getSys() string {
	sys := runtime.GOOS
	logrus.Infof("cur sys: %s", sys)
	return sys
}

// SendCommand
//
//	@Description: 执行命令
//	@receiver cli 需要执行命令的LocalCli对象
//	@param command 需要执行命令
//	@return string 命令结果
//	@return error 执行出错
func (cli *LocalCli) SendCommand(command string) (string, error) {
	if cli.cmd == nil || cli.in == nil || cli.standReader == nil || cli.errReader == nil {
		return "", errors.New("cmd not init")
	}
	logrus.Infof("send command: %s", command)
	_, err := cli.in.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}
	res, err := cli.GetResult()
	return res, err
}

func (cli *LocalCli) GetResult() (string, error) {
	for !cli.isCommandFinish() {
	}
	res := cli.standRes.String()
	errorRes := cli.errorRes.String()
	if errorRes != "" {
		logrus.Errorf("send command error: %s", errorRes)
		cli.errorRes.Reset()
		return "", errors.New(errorRes)
	}
	cli.standRes.Reset()
	return res, nil
}

func (cli *LocalCli) readOutput(reader *bufio.Reader, res *strings.Builder) {
	outputBytes := make([]byte, 1024)
	logrus.Infof("start read output")
	for {
		n, err := reader.Read(outputBytes) //获取屏幕的实时输出(并不是按照回车分割，所以要结合sumOutput)
		logrus.Infof("read output n: %d", n)
		if err != nil {
			if err == io.EOF {
				logrus.Infof("output over")
				break
			}
			logrus.Errorf("get output error: %+v", err)
		}
		if n > 0 {
			output := string(outputBytes[:n])
			res.WriteString(output)
			println(res.String())
		}
	}
}

func (cli *LocalCli) isCommandFinish() bool {
	if cli.isFinish(cli.standRes.String()) || cli.isFinish(cli.errorRes.String()) {
		return true
	}
	return false
}

func (cli *LocalCli) isFinish(output string) bool {
	if len(cli.overSign) == 0 {
		if cli.sys == "windows" {
			return cli.isFinishOnSign(output, WindowsOverSign)
		}
		return cli.isFinishOnSign(output, LinuxOverSign)
	}
	return cli.isFinishOnSign(output, cli.overSign)
}

func (cli *LocalCli) isFinishOnSign(output string, overSign []string) bool {
	for _, value := range overSign {
		if strings.HasSuffix(strings.TrimSpace(output), value) {
			return true
		}
	}
	return false
}
