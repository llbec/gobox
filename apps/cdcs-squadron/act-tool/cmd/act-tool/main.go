package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gobox/apps/cdcs-squadron/act-tool/internal/app"
	"gobox/apps/cdcs-squadron/act-tool/internal/config"
	actlib "gobox/tools/act"
)

var (
	configPath  = flag.String("config", "config.yaml", "配置文件路径")
	mode        = flag.String("mode", "interactive", "运行模式: interactive 或 batch")
	dryRun      = flag.Bool("dry-run", false, "仅模拟操作，不实际执行")
	interactive = flag.Bool("interactive", false, "操作前要求用户确认")
	targetPath  = flag.String("target", "", "覆盖目标路径")
	backupPath  = flag.String("backup", "", "覆盖备份路径")
	cmd         = flag.String("cmd", "", "批处理模式下的命令: archive/list/resend/modify/rollback")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 命令行参数覆盖配置
	if *targetPath != "" {
		cfg.TargetPath = *targetPath
	}
	if *backupPath != "" {
		cfg.BackupPath = *backupPath
	}

	cfg.DryRun = *dryRun
	cfg.Interactive = *interactive

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 创建 ActManager
	manager := actlib.NewActManager(cfg.TargetPath, cfg.BackupPath)

	// 创建应用实例
	application := app.NewApp(manager, cfg)

	// 根据模式选择运行方式
	var runErr error
	switch *mode {
	case "interactive":
		runErr = application.RunInteractive(ctx)
	case "batch":
		if *cmd == "" {
			fmt.Println("批处理模式下必须指定 --cmd 参数")
			os.Exit(1)
		}
		runErr = application.RunBatch(ctx, *cmd)
	default:
		fmt.Printf("未知的运行模式: %s\n", *mode)
		os.Exit(1)
	}

	if runErr != nil {
		log.Fatalf("执行失败: %v", runErr)
	}
}
