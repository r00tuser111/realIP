package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/realip/internal/logger"
	"github.com/yourusername/realip/internal/output"
	"github.com/yourusername/realip/internal/resolver"
	"github.com/yourusername/realip/pkg/banner"
)

// 命令行参数
type Config struct {
	Target     string // 单个目标域名
	ListFile   string // 域名列表文件
	JSONOutput bool   // 是否输出JSON格式
	CSVOutput  bool   // 是否输出CSV格式
	Threads    int    // 并发线程数
	OutputFile string // 输出文件路径
	Silent     bool   // 是否静默模式
	FilterCNAME bool  // 是否过滤CNAME记录
}

func main() {
	// 显示banner
	banner.Show()

	// 解析命令行参数
	config := parseFlags()

	// 初始化日志器
	log := logger.New(!config.Silent)

	// 检查参数有效性
	if config.Target == "" && config.ListFile == "" {
		log.Error("必须指定目标域名 (-u) 或域名列表文件 (-l)")
		flag.Usage()
		os.Exit(1)
	}

	// 如果同时指定了JSON和CSV，输出提示并默认使用JSON
	if config.JSONOutput && config.CSVOutput {
		log.Warn("同时指定了JSON和CSV格式，将默认使用JSON格式")
		config.CSVOutput = false
	}

	// 初始化解析器
	r := resolver.New(config.Threads, log, config.FilterCNAME)

	// 初始化输出处理器
	var out output.Writer
	if config.JSONOutput {
		out = output.NewJSONWriter()
	} else if config.CSVOutput {
		out = output.NewCSVWriter()
	} else {
		out = output.NewTextWriter()
	}

	// 设置输出位置
	var outputDest *os.File
	var err error
	if config.OutputFile != "" {
		outputDest, err = os.Create(config.OutputFile)
		if err != nil {
			log.Error("无法创建输出文件: %v", err)
			os.Exit(1)
		}
		defer outputDest.Close()
	} else {
		outputDest = os.Stdout
	}

	// 处理单个域名
	if config.Target != "" {
		log.Info("正在解析域名: %s", config.Target)
		result := r.ResolveSingle(config.Target)
		err = out.Write(outputDest, []resolver.Result{result})
		if err != nil {
			log.Error("输出结果时出错: %v", err)
			os.Exit(1)
		}
		log.Info("解析完成")
		return
	}

	// 处理域名列表
	if config.ListFile != "" {
		log.Info("正在从文件加载域名列表: %s", config.ListFile)
		domains, err := loadDomainsFromFile(config.ListFile)
		if err != nil {
			log.Error("加载域名列表失败: %v", err)
			os.Exit(1)
		}

		log.Info("开始解析 %d 个域名，使用 %d 个线程", len(domains), config.Threads)
		results := r.ResolveMultiple(domains)
		err = out.Write(outputDest, results)
		if err != nil {
			log.Error("输出结果时出错: %v", err)
			os.Exit(1)
		}
		log.Info("所有域名解析完成")
	}
}

// 解析命令行参数
func parseFlags() Config {
	var config Config

	// 定义命令行参数
	flag.StringVar(&config.Target, "target", "", "单个目标域名")
	flag.StringVar(&config.Target, "u", "", "单个目标域名 (简写)")
	
	flag.StringVar(&config.ListFile, "list", "", "域名列表文件路径 (每行一个域名)")
	flag.StringVar(&config.ListFile, "l", "", "域名列表文件路径 (简写)")
	
	flag.BoolVar(&config.JSONOutput, "json", false, "以JSON格式输出结果")
	flag.BoolVar(&config.CSVOutput, "csv", false, "以CSV格式输出结果")
	
	flag.IntVar(&config.Threads, "threads", 10, "并发线程数")
	flag.IntVar(&config.Threads, "t", 10, "并发线程数 (简写)")
	
	flag.StringVar(&config.OutputFile, "o", "", "输出文件路径")
	
	flag.BoolVar(&config.Silent, "silent", false, "关闭所有日志输出")
	
	flag.BoolVar(&config.FilterCNAME, "cname", false, "过滤CNAME记录，只显示直接解析到IP的域名")

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "选项:")
		flag.PrintDefaults()
	}

	flag.Parse()
	return config
}

// 从文件加载域名列表
func loadDomainsFromFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := []string{}
	for _, line := range splitLines(string(content)) {
		line = cleanLine(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

// 按行分割文本
func splitLines(s string) []string {
	var lines []string
	var line string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(r)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

// 清理行内容（去除空白字符等）
func cleanLine(s string) string {
	var result string
	for _, r := range s {
		if r != '\r' && r != '\t' && r != ' ' {
			result += string(r)
		}
	}
	return result
} 