package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leoobai/aigc-check/internal/analyzer"
	"github.com/leoobai/aigc-check/internal/config"
	"github.com/leoobai/aigc-check/internal/models"
	"github.com/leoobai/aigc-check/internal/reporter"
)

const version = "0.1.0"

func main() {
	// 定义命令行参数
	var (
		inputFile  string
		outputFile string
		format     string
		configFile string
		showHelp   bool
		showVer    bool
	)

	flag.StringVar(&inputFile, "f", "", "输入文件路径")
	flag.StringVar(&inputFile, "file", "", "输入文件路径")
	flag.StringVar(&outputFile, "o", "", "输出文件路径（可选）")
	flag.StringVar(&outputFile, "output", "", "输出文件路径（可选）")
	flag.StringVar(&format, "format", "text", "输出格式: text, json")
	flag.StringVar(&configFile, "c", "", "配置文件路径（可选）")
	flag.StringVar(&configFile, "config", "", "配置文件路径（可选）")
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")
	flag.BoolVar(&showVer, "v", false, "显示版本信息")
	flag.BoolVar(&showVer, "version", false, "显示版本信息")

	flag.Parse()

	// 显示帮助信息
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// 显示版本信息
	if showVer {
		fmt.Printf("AIGC-Check v%s\n", version)
		os.Exit(0)
	}

	// 检查输入文件
	if inputFile == "" {
		fmt.Fprintln(os.Stderr, "错误: 必须指定输入文件")
		fmt.Fprintln(os.Stderr, "使用 -h 或 --help 查看帮助信息")
		os.Exit(1)
	}

	// 运行检测
	if err := run(inputFile, outputFile, format, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 执行检测流程
func run(inputFile, outputFile, format, configFile string) error {
	// 加载配置
	var cfg *config.Config
	var err error

	if configFile != "" {
		cfg, err = config.LoadConfig(configFile)
		if err != nil {
			return fmt.Errorf("加载配置文件失败: %w", err)
		}
	} else {
		// 尝试加载默认配置
		defaultConfigPath := filepath.Join("configs", "aigc-check.yaml")
		cfg, err = config.LoadConfig(defaultConfigPath)
		if err != nil {
			// 使用默认配置
			defaultCfg := config.DefaultConfig
			cfg = &defaultCfg
		}
	}

	// 如果命令行指定了格式，覆盖配置
	if format != "" {
		cfg.Output.DefaultFormat = format
	}

	// 读取输入文件
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}

	text := string(content)
	if text == "" {
		return fmt.Errorf("输入文件为空")
	}

	// 创建分析器
	a := analyzer.NewAnalyzer(cfg)

	// 执行分析
	request := models.DetectionRequest{
		Text: text,
		Options: models.DetectionOptions{
			Language:     cfg.Output.Language,
			OutputFormat: cfg.Output.DefaultFormat,
		},
	}

	result, err := a.Analyze(request)
	if err != nil {
		return fmt.Errorf("分析失败: %w", err)
	}

	// 生成报告
	var rep reporter.Reporter
	switch cfg.Output.DefaultFormat {
	case "json":
		rep = reporter.NewJSONReporter(true)
	case "text":
		rep = reporter.NewTextReporter(cfg.Output.ColorEnabled)
	default:
		rep = reporter.NewTextReporter(cfg.Output.ColorEnabled)
	}

	report, err := rep.Generate(result)
	if err != nil {
		return fmt.Errorf("生成报告失败: %w", err)
	}

	// 输出报告
	if outputFile != "" {
		// 写入文件
		if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		fmt.Printf("报告已保存到: %s\n", outputFile)
	} else {
		// 输出到标准输出
		fmt.Println(report)
	}

	return nil
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("AIGC-Check - AI生成内容检测工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  aigc-check -f <文件路径> [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -f, --file <路径>      输入文件路径（必需）")
	fmt.Println("  -o, --output <路径>    输出文件路径（可选，默认输出到标准输出）")
	fmt.Println("  -format <格式>         输出格式: text, json（默认: text）")
	fmt.Println("  -c, --config <路径>    配置文件路径（可选）")
	fmt.Println("  -h, --help             显示帮助信息")
	fmt.Println("  -v, --version          显示版本信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 检测文本文件")
	fmt.Println("  aigc-check -f sample.txt")
	fmt.Println()
	fmt.Println("  # 使用JSON格式输出")
	fmt.Println("  aigc-check -f sample.txt -format json")
	fmt.Println()
	fmt.Println("  # 保存报告到文件")
	fmt.Println("  aigc-check -f sample.txt -o report.txt")
	fmt.Println()
	fmt.Println("  # 使用自定义配置")
	fmt.Println("  aigc-check -f sample.txt -c custom-config.yaml")
	fmt.Println()
}
