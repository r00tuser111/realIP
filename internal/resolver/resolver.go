package resolver

import (
	"net"
	"sync"
	"time"
	"fmt"

	"github.com/yourusername/realip/internal/logger"
)

// Result 表示域名解析结果
type Result struct {
	Domain    string `json:"domain"`    // 域名
	IP        string `json:"ip"`        // 解析到的IP
	Timestamp string `json:"timestamp"` // 解析时间戳
	Error     string `json:"error"`     // 错误信息（如果有）
	HasCNAME  bool   `json:"has_cname"` // 是否有CNAME记录
}

// Resolver 域名解析器
type Resolver struct {
	threads     int            // 并发线程数
	logger      *logger.Logger // 日志记录器
	filterCNAME bool          // 是否过滤CNAME记录
}

// New 创建新的解析器
func New(threads int, logger *logger.Logger, filterCNAME bool) *Resolver {
	return &Resolver{
		threads:     threads,
		logger:      logger,
		filterCNAME: filterCNAME,
	}
}

// ResolveSingle 解析单个域名
func (r *Resolver) ResolveSingle(domain string) Result {
	r.logger.Debug("解析域名: %s", domain)
	
	ip, hasCNAME, err := r.resolve(domain)
	
	result := Result{
		Domain:    domain,
		Timestamp: time.Now().Format(time.RFC3339),
		HasCNAME:  hasCNAME,
	}
	
	if err != nil {
		r.logger.Debug("解析域名 %s 失败: %v", domain, err)
		result.Error = err.Error()
	} else {
		r.logger.Debug("域名 %s 解析结果: %s (CNAME: %v)", domain, ip, hasCNAME)
		result.IP = ip
	}
	
	return result
}

// ResolveMultiple 批量解析多个域名
func (r *Resolver) ResolveMultiple(domains []string) []Result {
	r.logger.Debug("开始批量解析 %d 个域名", len(domains))
	
	resultsChan := make(chan Result, len(domains))
	var wg sync.WaitGroup
	
	// 限制并发数的通道
	semaphore := make(chan struct{}, r.threads)
	
	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			
			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			// 解析域名
			resultsChan <- r.ResolveSingle(d)
		}(domain)
	}
	
	// 等待所有goroutine完成
	wg.Wait()
	close(resultsChan)
	
	// 收集结果
	var results []Result
	for result := range resultsChan {
		// 如果启用了CNAME过滤，且域名有CNAME记录，则跳过该结果
		if r.filterCNAME && result.HasCNAME {
			r.logger.Debug("跳过带有CNAME记录的域名: %s", result.Domain)
			continue
		}
		results = append(results, result)
	}
	
	r.logger.Debug("批量解析完成，共处理 %d 个域名", len(results))
	return results
}

// resolve 实际的域名解析逻辑
func (r *Resolver) resolve(domain string) (string, bool, error) {
	// 首先检查是否有CNAME记录
	cname, err := net.LookupCNAME(domain)
	hasCNAME := err == nil && cname != domain+"."

	// 如果启用了CNAME过滤且存在CNAME记录，直接返回错误
	if r.filterCNAME && hasCNAME {
		return "", hasCNAME, fmt.Errorf("域名 %s 存在CNAME记录: %s", domain, cname)
	}

	// 获取所有IP地址
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", hasCNAME, err
	}

	// 检查是否只有一个IP地址
	if len(ips) != 1 {
		return "", hasCNAME, fmt.Errorf("域名 %s 有多个IP地址或无IP地址", domain)
	}

	// 返回唯一的IP地址
	if ipv4 := ips[0].To4(); ipv4 != nil {
		return ipv4.String(), hasCNAME, nil
	}
	return ips[0].String(), hasCNAME, nil
}