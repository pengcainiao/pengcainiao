package tests

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"
)

func TestRedis(T *testing.T) {
	//env.RedisAddr = "127.0.0.1:6379"
	//env.RedisMasterName = "119.29.5.54"
	//env.DbDSN = "penglonghui:Penglonghui!123!@tcp(119.29.5.54:3306)/okr?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	sshHost := "119.29.5.54:22"   // SSH服务器地址及端口
	username := "root"            // SSH用户名
	redisHost := "127.0.0.1:6379" // Redis服务器地址及端口
	// 私钥文件路径
	privateKeyPath := "id_ed25519"
	localPort := "6380"

	// 读取私钥文件
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// 解析私钥文件
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// 创建 SSH 客户端配置
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 这会忽略主机密钥验证，不建议在生产环境中使用
	}
	// 建立SSH连接
	sshClient, err := ssh.Dial("tcp", sshHost, sshConfig)
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	defer sshClient.Close()

	// 创建本地listener
	listener, err := net.Listen("tcp", "localhost:"+localPort)
	if err != nil {
		log.Fatalf("failed to listen on local port: %s", err)
	}
	defer listener.Close()

	// 在新goroutine中处理隧道连接
	go func() {
		for {
			localConn, err := listener.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %s", err)
				continue
			}

			remoteConn, err := sshClient.Dial("tcp", redisHost)
			if err != nil {
				log.Printf("failed to connect to remote host: %s", err)
				localConn.Close()
				continue
			}

			// 将数据从localConn复制到remoteConn
			go func() {
				defer localConn.Close()
				defer remoteConn.Close()
				_, err := io.Copy(remoteConn, localConn)
				if err != nil {
					log.Printf("io.Copy error: %s", err)
				}
			}()

			// 将数据从remoteConn复制到localConn
			go func() {
				defer localConn.Close()
				defer remoteConn.Close()
				_, err := io.Copy(localConn, remoteConn)
				if err != nil {
					log.Printf("io.Copy error: %s", err)
				}
			}()
		}
	}()

	// 等待一段时间让隧道建立起来
	time.Sleep(2 * time.Second)

	// 使用Go-Redis客户端连接到Redis服务器
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:" + localPort,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val := rdb.Get(ctx, "AAA").Val()
	fmt.Println("val:", val)

	keys := fmt.Sprintf("coin_agent_login_%v", "root")
	if rdb.Get(ctx, keys).Val() == "0" {
		fmt.Println("fail :", rdb.Get(ctx, keys).Val())
		return
	} else {
		fmt.Println("suc :", rdb.Get(ctx, keys).Val())
	}
}
