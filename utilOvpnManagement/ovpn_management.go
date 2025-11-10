package utilOvpnManagement

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Client 表示一个在线用户
type Client struct {
	ClientId         int    `json:"client_id"`
	CommonName       string `json:"common_name"`
	RealAddress      string `json:"real_address"`
	VirtualAddress   string `json:"virtual_address"`
	VirtualV6Address string `json:"virtual_v6_address"`
	BytesReceived    int64  `json:"bytes_received"`
	BytesSent        int64  `json:"bytes_sent"`
	ConnectedSince   string `json:"connected_since"`
	Username         string `json:"username"`
}

// Route 表示一个路由表项
type Route struct {
	VirtualAddress string `json:"virtual_address"`
	CommonName     string `json:"common_name"`
	RealAddress    string `json:"real_address"`
	LastRef        string `json:"last_ref"`
}

// StatusInfo 封装状态信息
type StatusInfo struct {
	NodeId   string   `json:"node_id,omitempty"`
	Title   string   `json:"title"`
	Time    string   `json:"time"`
	Clients []Client `json:"clients"`
	Routes  []Route  `json:"routes"`
}

type OpenVpnManagement struct {
	Addr     string
	Password string
	Timeout  time.Duration
	conn     net.Conn
	reader   *bufio.Reader
}

// NewOpenVpnManagement 创建新的 OpenVPN 管理客户端
func NewOpenVpnManagement(addr, password string, timeOut time.Duration) *OpenVpnManagement {
	return &OpenVpnManagement{
		Addr:     addr,
		Password: password,
		Timeout:  timeOut,
	}
}

// Connect 建立连接
func (m *OpenVpnManagement) Connect() (err error) {
	if nil != m.conn {
		return
	}
	conn, err := net.DialTimeout("tcp", m.Addr, m.Timeout)
	if err != nil {
		err = fmt.Errorf("连接失败: %v", err)
		return
	}
	m.conn = conn
	m.reader = bufio.NewReader(conn)

	// 如果需要密码
	if m.Password != "" {
		line, _ := m.reader.ReadString(':') // 读取 "Enter Management Password:"
		if !strings.Contains(strings.ToUpper(line), "PASSWORD") {
			err = fmt.Errorf("未收到密码提示")
			return
		}
		m.conn.Write([]byte(m.Password + "\n"))
	}

	return nil
}

// Close 关闭连接
func (m *OpenVpnManagement) Close() {
	if m.conn != nil {
		m.conn.Write([]byte("exit\n"))
		m.conn.Close()
		m.conn = nil
	}
}

// RunCommand 执行任意命令并返回原始输出
func (m *OpenVpnManagement) RunCommand(cmd string, endStr string) (result string, err error) {
	err = m.Connect()
	if nil != err {
		return
	}
	defer m.Close()

	if m.conn == nil {
		return "", fmt.Errorf("未连接")
	}
	m.conn.Write([]byte(cmd + "\n"))
	var output strings.Builder

	for {
		line, err := m.reader.ReadString('\n')
		if err != nil {
			break
		}
		output.WriteString(line)

		if "" == endStr || strings.Contains(line, endStr) {
			break
		}
	}
	result = output.String()
	return
}

// GetStatus 解析 status 2 输出
func (m *OpenVpnManagement) GetStatus() (info *StatusInfo, err error) {
	raw, err := m.RunCommand("status 2", "END")
	if err != nil {
		return
	}
	lines := strings.Split(raw, "\n")
	info = &StatusInfo{}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		switch parts[0] {
		case "TITLE":
			info.Title = parts[1]
		case "TIME":
			if len(parts) > 1 {
				info.Time = parts[1]
			}
		case "CLIENT_LIST":
			if len(parts) < 13 {
				continue
			}
			clientId, _ := strconv.Atoi(parts[10])
			client := Client{
				ClientId:         clientId,
				CommonName:       parts[1],
				RealAddress:      parts[2],
				VirtualAddress:   parts[3],
				VirtualV6Address: parts[4],
				ConnectedSince:   parts[7],
				Username:         parts[9],
			}
			client.BytesReceived, _ = strconv.ParseInt(parts[5], 10, 64)
			client.BytesSent, _ = strconv.ParseInt(parts[6], 10, 64)
			info.Clients = append(info.Clients, client)
		case "ROUTING_TABLE":
			if len(parts) < 5 {
				continue
			}
			info.Routes = append(info.Routes, Route{
				VirtualAddress: parts[1],
				CommonName:     parts[2],
				RealAddress:    parts[3],
				LastRef:        parts[4],
			})
		}
	}
	return
}

// KickClient 根据 client id 断开连接
func (m *OpenVpnManagement) KickClient(clientId string) (err error) {
	cmd := fmt.Sprintf("client-kill %s", clientId)
	_, err = m.RunCommand(cmd, "client-kill command")
	return err
}

// KickClientByCommonName 根据 CommonName 断开连接
func (m *OpenVpnManagement) KickClientByCommonName(name string) (err error) {
	cmd := fmt.Sprintf("kill %s", name)
	_, err = m.RunCommand(cmd, "")
	return err
}
