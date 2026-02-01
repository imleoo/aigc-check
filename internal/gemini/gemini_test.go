package gemini

import (
	"context"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Model != "gemini-pro" {
		t.Errorf("Expected model gemini-pro, got %s", cfg.Model)
	}

	if cfg.Temperature != 0.3 {
		t.Errorf("Expected temperature 0.3, got %f", cfg.Temperature)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", cfg.Timeout)
	}

	if !cfg.Cache.Enabled {
		t.Error("Expected cache to be enabled by default")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "未启用时不需要验证",
			config: Config{
				Enabled: false,
				APIKey:  "",
			},
			wantError: false,
		},
		{
			name: "启用但无API Key",
			config: Config{
				Enabled: true,
				APIKey:  "",
			},
			wantError: true,
		},
		{
			name: "有效配置",
			config: Config{
				Enabled:     true,
				APIKey:      "test-api-key",
				Model:       "gemini-pro",
				Temperature: 0.3,
				MaxTokens:   500,
				Timeout:     30 * time.Second,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name: "启用且有API Key",
			config: Config{
				Enabled: true,
				APIKey:  "test-key",
			},
			expected: true,
		},
		{
			name: "启用但无API Key",
			config: Config{
				Enabled: true,
				APIKey:  "",
			},
			expected: false,
		},
		{
			name: "未启用但有API Key",
			config: Config{
				Enabled: false,
				APIKey:  "test-key",
			},
			expected: false,
		},
		{
			name: "未启用且无API Key",
			config: Config{
				Enabled: false,
				APIKey:  "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.IsEnabled(); got != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCache(t *testing.T) {
	cfg := CacheConfig{
		Enabled:    true,
		TTL:        1 * time.Hour,
		MaxEntries: 10,
	}

	cache := NewCache(cfg)

	// 测试基本的 Get/Set
	t.Run("基本操作", func(t *testing.T) {
		cache.Set("key1", "value1")
		value, found := cache.Get("key1")

		if !found {
			t.Error("Expected to find cached value")
		}
		if value != "value1" {
			t.Errorf("Expected value1, got %s", value)
		}
	})

	// 测试不存在的 key
	t.Run("不存在的key", func(t *testing.T) {
		_, found := cache.Get("nonexistent")
		if found {
			t.Error("Expected not to find nonexistent key")
		}
	})

	// 测试删除
	t.Run("删除", func(t *testing.T) {
		cache.Set("key2", "value2")
		cache.Delete("key2")
		_, found := cache.Get("key2")
		if found {
			t.Error("Expected key2 to be deleted")
		}
	})

	// 测试清空
	t.Run("清空", func(t *testing.T) {
		cache.Set("key3", "value3")
		cache.Set("key4", "value4")
		cache.Clear()
		if cache.Size() != 0 {
			t.Errorf("Expected cache size 0, got %d", cache.Size())
		}
	})

	// 测试最大条目数限制
	t.Run("最大条目数", func(t *testing.T) {
		cache.Clear()
		for i := 0; i < 15; i++ {
			cache.Set(string(rune('a'+i)), "value")
		}
		if cache.Size() > cfg.MaxEntries {
			t.Errorf("Cache size %d exceeds max %d", cache.Size(), cfg.MaxEntries)
		}
	})
}

func TestCache_Disabled(t *testing.T) {
	cfg := CacheConfig{
		Enabled: false,
	}

	cache := NewCache(cfg)
	cache.Set("key", "value")
	_, found := cache.Get("key")

	if found {
		t.Error("Expected cache to be disabled")
	}
}

func TestNewClient_Disabled(t *testing.T) {
	cfg := Config{
		Enabled: false,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Errorf("Expected no error for disabled config, got %v", err)
	}

	if client == nil {
		t.Error("Expected client to be created even when disabled")
	}
}

func TestClient_GenerateContent_NotEnabled(t *testing.T) {
	cfg := Config{
		Enabled: false,
	}

	client, _ := NewClient(cfg)
	_, err := client.GenerateContent(context.Background(), "test prompt")

	if err != ErrNotEnabled {
		t.Errorf("Expected ErrNotEnabled, got %v", err)
	}
}

func TestParseJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name:  "纯JSON",
			input: `{"key": "value"}`,
			expected: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:  "带代码块的JSON",
			input: "```json\n{\"key\": \"value\"}\n```",
			expected: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:  "带前导文本的JSON",
			input: "Here is the result:\n{\"key\": \"value\"}",
			expected: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			input:   "not json at all",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := parseJSONResponse(tt.input, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSONResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result["key"] != tt.expected["key"] {
				t.Errorf("Expected key=%v, got %v", tt.expected["key"], result["key"])
			}
		})
	}
}

func TestAnalyzer_New(t *testing.T) {
	cfg := Config{
		Enabled: false,
	}
	client, _ := NewClient(cfg)
	analyzer := NewAnalyzer(client)

	if analyzer == nil {
		t.Error("Expected analyzer to be created")
	}
}

func TestSuggester_New(t *testing.T) {
	cfg := Config{
		Enabled: false,
	}
	client, _ := NewClient(cfg)
	suggester := NewSuggester(client)

	if suggester == nil {
		t.Error("Expected suggester to be created")
	}
}

func TestSuggester_DefaultSuggestions(t *testing.T) {
	cfg := Config{
		Enabled: false,
	}
	client, _ := NewClient(cfg)
	suggester := NewSuggester(client)

	issues := []string{
		"检测到高频AI词汇",
		"句式结构单一",
		"缺乏个人化表达",
	}

	suggestions := suggester.getDefaultSuggestions(issues)

	if len(suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(suggestions))
	}

	for i, s := range suggestions {
		if s.Priority != i+1 {
			t.Errorf("Expected priority %d, got %d", i+1, s.Priority)
		}
	}
}
