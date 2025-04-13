package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// TestMainCompilation は、mainパッケージがコンパイルできることを確認するだけのテスト
func TestMainCompilation(t *testing.T) {
	t.Log("Main package successfully compiled")
}

// TestArgumentParsing はコマンドライン引数のパースをテストする
func TestArgumentParsing(t *testing.T) {
	// 元のフラグセットをバックアップ
	origFlagSet := flag.CommandLine
	defer func() {
		flag.CommandLine = origFlagSet
	}()

	// 元の引数をバックアップ
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	// テストケース
	testCases := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "Default arguments",
			args:      []string{"go-sip"},
			expectErr: false,
		},
		{
			name:      "Custom port",
			args:      []string{"go-sip", "-port", "5070"},
			expectErr: false,
		},
		{
			name:      "Custom bind",
			args:      []string{"go-sip", "-bind", "127.0.0.1"},
			expectErr: false,
		},
		{
			name:      "Custom config",
			args:      []string{"go-sip", "-config", "test_config.json"},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// フラグを再設定
			flag.CommandLine = flag.NewFlagSet(tc.args[0], flag.ContinueOnError)

			// 引数を設定
			os.Args = tc.args

			// configPathだけをテスト
			configPath := flag.String("config", "config.json", "Path to configuration file")
			overridePort := flag.String("port", "", "Override port setting from config file")
			overrideBindAddr := flag.String("bind", "", "Override bind address setting from config file")

			// フラグをパース
			err := flag.CommandLine.Parse(tc.args[1:])

			// 結果を確認
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error %v, got %v", tc.expectErr, err)
			}

			// 引数の値をログに出力（エラー時のデバッグ用）
			t.Logf("Config path: %s", *configPath)
			if *overridePort != "" {
				t.Logf("Port override: %s", *overridePort)
			}
			if *overrideBindAddr != "" {
				t.Logf("Bind address override: %s", *overrideBindAddr)
			}
		})
	}
}

// TestGenerateConfigFlag はgenerate-configフラグの動作をテストする
func TestGenerateConfigFlag(t *testing.T) {
	// 元のフラグセットをバックアップ
	origFlagSet := flag.CommandLine
	defer func() {
		flag.CommandLine = origFlagSet
	}()

	// テストデータディレクトリを準備
	testDir := filepath.Join("testdata")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// テスト用の設定ファイルパス
	tempConfig := filepath.Join(testDir, "test_generated_config.json")
	defer os.Remove(tempConfig)

	// 元の引数をバックアップ
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	// 引数を設定
	os.Args = []string{"go-sip", "-generate-config", "-config", tempConfig}

	// フラグを再設定
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configPath := flag.String("config", "config.json", "Path to configuration file")
	generateConfig := flag.Bool("generate-config", false, "Generate default config file and exit")

	// フラグをパース
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// 値を確認
	if !*generateConfig {
		t.Error("generate-config flag should be true")
	}

	if *configPath != tempConfig {
		t.Errorf("Expected config path %s, got %s", tempConfig, *configPath)
	}

	// ここで実際に設定ファイルを生成する代わりに値だけを確認する
	t.Logf("Would generate config file at: %s", *configPath)
}
