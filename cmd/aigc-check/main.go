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

const version = "2.0.0"

// runOptions 运行选项
type runOptions struct {
	inputFile        string
	outputFile       string
	format           string
	configFile       string
	enableMultimodal bool
	enableStatistics bool
	enableGemini     bool
	geminiAPIKey     string
	verbose          bool
}

func main() {
	// 定义命令行参数
	var (
		inputFile       string
		outputFile      string
		format          string
		configFile      string
		showHelp        bool
		showVer         bool
		enableMultimodal bool
		enableStatistics bool
		enableGemini    bool
		geminiAPIKey    string
		verbose         bool
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

	// 多模态检测参数
	flag.BoolVar(&enableMultimodal, "m", false, "启用多模态检测")
	flag.BoolVar(&enableMultimodal, "multimodal", false, "启用多模态检测")
	flag.BoolVar(&enableStatistics, "s", false, "启用统计分析层")
	flag.BoolVar(&enableStatistics, "statistics", false, "启用统计分析层")
	flag.BoolVar(&enableGemini, "g", false, "启用 Gemini API 语义分析")
	flag.BoolVar(&enableGemini, "gemini", false, "启用 Gemini API 语义分析")
	flag.StringVar(&geminiAPIKey, "api-key", "", "Gemini API Key（也可通过环境变量 GEMINI_API_KEY 设置）")
	flag.BoolVar(&verbose, "verbose", false, "显示详细分析结果")

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
	opts := runOptions{
		inputFile:        inputFile,
		outputFile:       outputFile,
		format:           format,
		configFile:       configFile,
		enableMultimodal: enableMultimodal,
		enableStatistics: enableStatistics,
		enableGemini:     enableGemini,
		geminiAPIKey:     geminiAPIKey,
		verbose:          verbose,
	}
	if err := run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 执行检测流程
func run(opts runOptions) error {
	// 加载配置
	var cfg *config.Config
	var err error

	if opts.configFile != "" {
		cfg, err = config.LoadConfig(opts.configFile)
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
	if opts.format != "" {
		cfg.Output.DefaultFormat = opts.format
	}

	// 应用多模态检测配置
	if opts.enableMultimodal {
		cfg.Multimodal.Enabled = true
	}
	if opts.enableStatistics {
		cfg.Multimodal.EnableStatistics = true
	}
	if opts.enableGemini {
		cfg.Multimodal.EnableSemantic = true
		cfg.Gemini.Enabled = true
	}
	if opts.geminiAPIKey != "" {
		cfg.Gemini.APIKey = opts.geminiAPIKey
	}
	if opts.verbose {
		cfg.Output.Verbose = true
	}

	// 读取输入文件
	content, err := os.ReadFile(opts.inputFile)
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
	if opts.outputFile != "" {
		// 写入文件
		if err := os.WriteFile(opts.outputFile, []byte(report), 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		fmt.Printf("报告已保存到: %s\n", opts.outputFile)
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
	fmt.Println("多模态检测选项:")
	fmt.Println("  -m, --multimodal       启用多模态检测（默认: false）")
	fmt.Println("  -s, --statistics       启用统计分析层（默认: false）")
	fmt.Println("  -g, --gemini           启用 Gemini API 语义分析（默认: false）")
	fmt.Println("  --api-key <key>        Gemini API Key（也可通过环境变量 GEMINI_API_KEY 设置）")
	fmt.Println("  --verbose              显示详细分析结果（默认: false）")
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
	fmt.Println("  # 启用多模态检测")
	fmt.Println("  aigc-check -f sample.txt -m")
	fmt.Println()
	fmt.Println("  # 启用统计分析层")
	fmt.Println("  aigc-check -f sample.txt -m -s")
	fmt.Println()
	fmt.Println("  # 启用完整多模态检测（包括 Gemini API）")
	fmt.Println("  aigc-check -f sample.txt -m -s -g --api-key YOUR_API_KEY")
	fmt.Println()
	fmt.Println("  # 显示详细分析结果")
	fmt.Println("  aigc-check -f sample.txt -m --verbose")
	fmt.Println()
}
