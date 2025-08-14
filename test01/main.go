package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

	"golang.org/x/crypto/ssh"
)

// go ssh 连接ssh
// 参考blog：
//
//	https://www.cnblogs.com/zhzhlong/p/12552410.html
//	https://blog.csdn.net/Naisu_kun/article/details/130598129
func main() {
	log.Println("main ...")
	client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("root")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH dial error: %s", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("new session error: %s", err.Error())
	}

	// session run执行命令
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("pwd"); err != nil {
		log.Fatalln("Failed to run: " + err.Error())
	}
	log.Println(strings.Trim(string(b.String()), "\n"))
	session.Close()

	// session执行Output命令
	session, _ = client.NewSession()
	result, err := session.Output("pwd")
	if err != nil {
		log.Fatalln("Failed to run command, Err:", err.Error())
	}
	log.Println(strings.Trim(string(result), "\n"))
	session.Close()

	// 模拟terminal
	session, _ = client.NewSession()
	// 会话输出关联到系统标准输出设备
	session.Stdout = os.Stdout
	// 会话错误输出关联到系统标准错误输出设备
	session.Stderr = os.Stderr
	// 会话输入关联到系统标准输入设备
	session.Stdin = os.Stdin

	// 设置终端参数
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}

	// 获取当前标准输出终端窗口尺寸 // 该操作可能有的平台上不可用，那么下面手动指定终端尺寸即可
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("unable to terminal.GetSize: ", err)
	}

	// 设置虚拟终端与远程会话关联
	if err = session.RequestPty("linux", termHeight, termWidth, modes); err != nil {
		log.Fatalf("request pty error: %s", err.Error())
	}

	// 启动远程Shell
	if err = session.Shell(); err != nil {
		log.Fatalf("start shell error: %s", err.Error())
	}
	// 启动远程Shell
	if err = session.Wait(); err != nil {
		log.Fatalf("return error: %s", err.Error())
	}
	log.Println("success ..")
}
