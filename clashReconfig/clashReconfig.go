package main

import (
	"bytes"
	"clashConfigOverwrite/common/clashConfig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var TemplateDir = ""
var defaultPrefix = ""

func main() {
	//读取主配置
	content, err := os.ReadFile("a.yaml")
	checkErr(err)
	mainData, err := clashConfig.Yaml2Struct(content)
	checkErr(err)
	//为默认节点加前缀
	mainData.Proxies = clashConfig.AddProxyNamePrefixSuffix(mainData.Proxies, GenPrefix(defaultPrefix), "")
	//获取节点提供者所有节点并加前缀(递归调用，应对第三方订阅也使用提供者情况)
	mainData.Proxies = append(mainData.Proxies, GetProxiesFromProviders(mainData.ProxyProviders)...)
	//列出所有自定义节点名称
	proxiesNames := Slice2Map(GetProxiesNames(mainData.Proxies))
	//列出所有组名称
	proxyGroupsNames := Slice2Map(GetProxyGroupNames(mainData.ProxyGroups))
	//保留代理名
	normalProxies := map[string]string{
		"direct": "direct", "reject": "reject",
	}
	//应用节点组use的提供者
	for _, group := range mainData.ProxyGroups {
		for i, proxyName := range group.Proxies {
			if _, ok := normalProxies[strings.ToLower(group.Proxies[i])]; ok {
				//direct reject
				group.Proxies[i] = group.Proxies[i]
			} else if _otherGroupName, ok := proxyGroupsNames[group.Proxies[i]]; ok && _otherGroupName != group.Name {
				//可以将其他组引入本组
				group.Proxies[i] = _otherGroupName
			} else {
				//自定义节点
				group.Proxies[i] = GenPrefix(defaultPrefix) + proxyName
				if _, ok := proxiesNames[group.Proxies[i]]; !ok {
					log.Fatalf("[warn] proxy-groups.%s.%s not found", group.Name, proxyName)
				}
			}
		}
		for i, useName := range group.Use {
			prefix := GenPrefix(useName)
			reg := regexp.MustCompile("^" + prefix)
			for _, proxy := range mainData.Proxies {
				proxyName := proxy["name"].(string)
				if reg.MatchString(proxyName) {
					group.Proxies = append(group.Proxies, proxyName)
				}
			}
			group.Use[i] = ""
		}
		group.Use = []string{}
	}
	mainData.ProxyProviders = nil
	output, err := yaml.Marshal(&mainData)
	checkErr(err)
	//获取规则提供者(无需递归，规则提供者只有规则)
	//按RULE-SET对应(节点组/节点)关系创建规则
	//输出新配置
	os.WriteFile("newConfig.yaml", escape2Utf8(output), 0777)
}
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetProxiesFromProviders(providers map[string]clashConfig.ProxyProvider) []clashConfig.Proxy {
	sumProxies := []clashConfig.Proxy{}
	var err error
	for name, provider := range providers {
		proxies := []clashConfig.Proxy{}
		content := []byte{}
		providerType := ""
		if provider.Url != "" {
			providerType = provider.Url
			content = HttpGet(provider.Url)
		} else {
			providerType = provider.Path
			content, err = os.ReadFile(path.Join(provider.Path))
			if err != nil {
				log.Println("[warn] ", err)
			}
		}
		if content == nil {
			continue
		}
		tempData, err := clashConfig.Yaml2Struct(content)
		if err != nil {
			log.Println("[error] provider error:", providerType, "\n", err)
			panic(err)
		}
		if len(tempData.ProxyProviders) > 0 {
			proxies = append(proxies, GetProxiesFromProviders(tempData.ProxyProviders)...)
		}
		proxies = append(proxies, tempData.Proxies...)
		proxies = clashConfig.AddProxyNamePrefixSuffix(proxies, name+"::", "")
		sumProxies = append(sumProxies, proxies...)
	}
	return sumProxies
}
func HttpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("warn", err)
		return []byte{}
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("warn", err)
		return []byte{}
	}
	return bts
}
func GenPrefix(str string) string {
	return str + "."
}
func GenSuffix(str string) string {
	return "." + str
}
func GetProxiesNames(data []clashConfig.Proxy) (result []string) {
	for _, item := range data {
		result = append(result, item["name"].(string))
	}
	return
}
func Slice2Map(data []string) map[string]string {
	res := map[string]string{}
	for _, str := range data {
		res[str] = str
	}
	return res
}
func GetProxyGroupNames(data []*clashConfig.ProxyGroup) []string {
	res := []string{}
	for _, group := range data {
		if group.Name == "" {
			continue
		}
		res = append(res, group.Name)
	}
	return res
}
func escape2Utf8(data []byte) []byte {
	re := regexp.MustCompile(`(\\U[0-9a-fA-F]{8})+`)
	for _, match := range re.FindAll(data, -1) {
		str, _ := strconv.Unquote(`"` + string(match) + `"`)
		data = bytes.ReplaceAll(data, match, []byte(str))
	}
	return data
}
func sliceDelStr(data []string, find string) []string {
	m := make(map[string]bool)
	for _, v := range data {
		if v != find {
			m[v] = true
		}
	}
	s2 := make([]string, 0, len(m))
	for k := range m {
		s2 = append(s2, k)
	}
	return s2
}
