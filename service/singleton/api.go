package singleton

import (
	"sync"

	"github.com/midoks/nezha/model"
	"github.com/midoks/nezha/pkg/utils"
)

var (
	ApiTokenList         = make(map[string]*model.ApiToken)
	UserIDToApiTokenList = make(map[uint64][]string)
	ApiLock              sync.RWMutex

	ServerAPI = &ServerAPIService{}
)

type ServerAPIService struct{}

// CommonResponse 常规返回结构 包含状态码 和 状态信息
type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CommonServerInfo struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	LastActive int64  `json:"last_active"`
	IPV4       string `json:"ipv4"`
	IPV6       string `json:"ipv6"`
	ValidIP    string `json:"valid_ip"`
}

// StatusResponse 服务器状态子结构 包含服务器信息与状态信息
type StatusResponse struct {
	CommonServerInfo
	Host   *model.Host      `json:"host"`
	Status *model.HostState `json:"status"`
}

// ServerStatusResponse 服务器状态返回结构 包含常规返回结构 和 服务器状态子结构
type ServerStatusResponse struct {
	CommonResponse
	Result []*StatusResponse `json:"result"`
}

// ServerInfoResponse 服务器信息返回结构 包含常规返回结构 和 服务器信息子结构
type ServerInfoResponse struct {
	CommonResponse
	Result []*CommonServerInfo `json:"result"`
}

func InitAPI() {
	ApiTokenList = make(map[string]*model.ApiToken)
	UserIDToApiTokenList = make(map[uint64][]string)
}

func LoadAPI() {
	InitAPI()
	var tokenList []*model.ApiToken
	DB.Find(&tokenList)
	for _, token := range tokenList {
		ApiTokenList[token.Token] = token
		UserIDToApiTokenList[token.UserID] = append(UserIDToApiTokenList[token.UserID], token.Token)
	}
}

// GetStatusByIDList 获取传入IDList的服务器状态信息
func (s *ServerAPIService) GetStatusByIDList(idList []uint64) *ServerStatusResponse {
	res := &ServerStatusResponse{}
	res.Result = make([]*StatusResponse, 0)

	ServerLock.RLock()
	defer ServerLock.RUnlock()

	for _, v := range idList {
		server := ServerList[v]
		if server == nil {
			continue
		}
		ipv4, ipv6, validIP := utils.SplitIPAddr(server.Host.IP)
		info := CommonServerInfo{
			ID:         server.ID,
			Name:       server.Name,
			Tag:        server.Tag,
			LastActive: server.LastActive.Unix(),
			IPV4:       ipv4,
			IPV6:       ipv6,
			ValidIP:    validIP,
		}
		res.Result = append(res.Result, &StatusResponse{
			CommonServerInfo: info,
			Host:             server.Host,
			Status:           server.State,
		})
	}
	res.CommonResponse = CommonResponse{
		Code:    0,
		Message: "success",
	}
	return res
}

// GetStatusByTag 获取传入分组的所有服务器状态信息
func (s *ServerAPIService) GetStatusByTag(tag string) *ServerStatusResponse {
	return s.GetStatusByIDList(ServerTagToIDList[tag])
}

// GetAllStatus 获取所有服务器状态信息
func (s *ServerAPIService) GetAllStatus() *ServerStatusResponse {
	res := &ServerStatusResponse{}
	res.Result = make([]*StatusResponse, 0)
	ServerLock.RLock()
	defer ServerLock.RUnlock()
	for _, v := range ServerList {
		host := v.Host
		state := v.State
		if host == nil || state == nil {
			continue
		}
		ipv4, ipv6, validIP := utils.SplitIPAddr(host.IP)
		info := CommonServerInfo{
			ID:         v.ID,
			Name:       v.Name,
			Tag:        v.Tag,
			LastActive: v.LastActive.Unix(),
			IPV4:       ipv4,
			IPV6:       ipv6,
			ValidIP:    validIP,
		}
		res.Result = append(res.Result, &StatusResponse{
			CommonServerInfo: info,
			Host:             v.Host,
			Status:           v.State,
		})
	}
	res.CommonResponse = CommonResponse{
		Code:    0,
		Message: "success",
	}
	return res
}

// GetListByTag 获取传入分组的所有服务器信息
func (s *ServerAPIService) GetListByTag(tag string) *ServerInfoResponse {
	res := &ServerInfoResponse{}
	res.Result = make([]*CommonServerInfo, 0)

	ServerLock.RLock()
	defer ServerLock.RUnlock()
	for _, v := range ServerTagToIDList[tag] {
		host := ServerList[v].Host
		if host == nil {
			continue
		}
		ipv4, ipv6, validIP := utils.SplitIPAddr(host.IP)
		info := &CommonServerInfo{
			ID:         v,
			Name:       ServerList[v].Name,
			Tag:        ServerList[v].Tag,
			LastActive: ServerList[v].LastActive.Unix(),
			IPV4:       ipv4,
			IPV6:       ipv6,
			ValidIP:    validIP,
		}
		res.Result = append(res.Result, info)
	}
	res.CommonResponse = CommonResponse{
		Code:    0,
		Message: "success",
	}
	return res
}

// GetAllList 获取所有服务器信息
func (s *ServerAPIService) GetAllList() *ServerInfoResponse {
	res := &ServerInfoResponse{}
	res.Result = make([]*CommonServerInfo, 0)

	ServerLock.RLock()
	defer ServerLock.RUnlock()
	for _, v := range ServerList {
		host := v.Host
		if host == nil {
			continue
		}
		ipv4, ipv6, validIP := utils.SplitIPAddr(host.IP)
		info := &CommonServerInfo{
			ID:         v.ID,
			Name:       v.Name,
			Tag:        v.Tag,
			LastActive: v.LastActive.Unix(),
			IPV4:       ipv4,
			IPV6:       ipv6,
			ValidIP:    validIP,
		}
		res.Result = append(res.Result, info)
	}
	res.CommonResponse = CommonResponse{
		Code:    0,
		Message: "success",
	}
	return res
}
