package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourusername/realip/internal/resolver"
)

// Writer 输出写入器接口
type Writer interface {
	Write(w io.Writer, results []resolver.Result) error
}

// JSONWriter JSON格式输出
type JSONWriter struct{}

// NewJSONWriter 创建JSON输出写入器
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

// Write 将结果以JSON格式写入
func (j *JSONWriter) Write(w io.Writer, results []resolver.Result) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

// CSVWriter CSV格式输出
type CSVWriter struct{}

// NewCSVWriter 创建CSV输出写入器
func NewCSVWriter() *CSVWriter {
	return &CSVWriter{}
}

// Write 将结果以CSV格式写入
func (c *CSVWriter) Write(w io.Writer, results []resolver.Result) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入CSV表头
	if err := writer.Write([]string{"Domain", "IP", "Has_CNAME", "Timestamp", "Error"}); err != nil {
		return err
	}

	// 写入数据行
	for _, result := range results {
		if err := writer.Write([]string{
			result.Domain,
			result.IP,
			fmt.Sprintf("%v", result.HasCNAME),
			result.Timestamp,
			result.Error,
		}); err != nil {
			return err
		}
	}

	return nil
}

// TextWriter 纯文本格式输出
type TextWriter struct{}

// NewTextWriter 创建文本输出写入器
func NewTextWriter() *TextWriter {
	return &TextWriter{}
}

// Write 将结果以纯文本格式写入
func (t *TextWriter) Write(w io.Writer, results []resolver.Result) error {
	// 计算最长域名，以便对齐输出
	maxDomainLen := 0
	for _, result := range results {
		if len(result.Domain) > maxDomainLen {
			maxDomainLen = len(result.Domain)
		}
	}
	
	// 打印表头
	fmt.Fprintf(w, "%-*s  %-15s  %-8s  %s\n", maxDomainLen, "域名", "IP地址", "CNAME", "状态")
	fmt.Fprintf(w, "%s\n", repeatChar('-', maxDomainLen+2+15+2+8+2+10))
	
	// 打印每行结果
	for _, result := range results {
		status := "成功"
		ip := result.IP
		cname := "否"
		
		if result.HasCNAME {
			cname = "是"
		}
		
		if result.Error != "" {
			status = "失败"
			ip = "-"
		}
		
		fmt.Fprintf(w, "%-*s  %-15s  %-8s  %s\n", maxDomainLen, result.Domain, ip, cname, status)
	}
	
	return nil
}

// repeatChar 重复字符串
func repeatChar(char byte, count int) string {
	result := make([]byte, count)
	for i := 0; i < count; i++ {
		result[i] = char
	}
	return string(result)
} 