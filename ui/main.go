package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	// WebSocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Store active analysis sessions
	sessions = make(map[string]*AnalysisSession)
	sessionMu sync.RWMutex
)

// AnalysisSession represents an active analysis
type AnalysisSession struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"` // running, completed, failed
	Output        []string               `json:"output"`
	Result        string                 `json:"result"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time,omitempty"`
	CompanyName   string                 `json:"company_name"`
	StockCode     string                 `json:"stock_code"`
	Market        string                 `json:"market"`
	subscribers   []*websocket.Conn
	mu            sync.RWMutex
}

// AnalysisRequest represents the analysis request from frontend
type AnalysisRequest struct {
	CompanyName string `json:"company_name"`
	StockCode   string `json:"stock_code"`
	Market      string `json:"market"`
}

// AISettings represents the AI model configuration
type AISettings struct {
	APIKey      string  `json:"api_key"`
	BaseURL     string  `json:"base_url"`
	ModelName   string  `json:"model_name"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// getSettingsFilePath returns the path to the settings file
func getSettingsFilePath() string {
	// Store settings in user's home directory or current directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "ai_settings.json"
	}
	settingsDir := filepath.Join(homeDir, ".a_stock_analysis")
	os.MkdirAll(settingsDir, 0755)
	return filepath.Join(settingsDir, "ai_settings.json")
}

// loadSettings loads AI settings from file
func loadSettings() (*AISettings, error) {
	settingsPath := getSettingsFilePath()
	
	// Check if file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// Return default settings
		return &AISettings{
			APIKey:      "",
			BaseURL:     "https://api.openai.com/v1",
			ModelName:   "gpt-4o",
			Temperature: 0.8,
			MaxTokens:   14000,
		}, nil
	}
	
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, err
	}
	
	var settings AISettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	
	return &settings, nil
}

// saveSettings saves AI settings to file
func saveSettings(settings *AISettings) error {
	settingsPath := getSettingsFilePath()
	
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(settingsPath, data, 0644)
}

// writeSettingsToEnv writes settings to .env file for Python to use
func writeSettingsToEnv(settings *AISettings) error {
	// Find the Python script directory
	cliEntryPath := findCliEntryPy()
	if cliEntryPath == "" {
		return fmt.Errorf("未找到Python分析脚本目录 (cli_entry.py)")
	}
	
	pythonDir := filepath.Dir(cliEntryPath)
	envPath := filepath.Join(pythonDir, ".env")
	
	// Create env content
	envContent := fmt.Sprintf(`OPENAI_API_KEY=%s
OPENAI_BASE_URL=%s
OPENAI_MODEL_NAME=%s
TEMPERATURE=%f
MAX_TOKENS=%d
`, settings.APIKey, settings.BaseURL, settings.ModelName, settings.Temperature, settings.MaxTokens)
	
	return os.WriteFile(envPath, []byte(envContent), 0644)
}

func getSettingsHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := loadSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request) {
	var settings AISettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Save to settings file
	if err := saveSettings(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Also write to .env for Python
	if err := writeSettingsToEnv(&settings); err != nil {
		log.Printf("警告: 无法写入.env文件: %v", err)
		// Don't fail the request, just log the warning
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "设置已保存",
	})
}

func main() {
	r := mux.NewRouter()

	// Serve static files
	staticDir := "./static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// API endpoints
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/settings", settingsPageHandler).Methods("GET")
	r.HandleFunc("/api/analyze", analyzeHandler).Methods("POST")
	r.HandleFunc("/api/session/{id}", getSessionHandler).Methods("GET")
	r.HandleFunc("/api/settings", getSettingsHandler).Methods("GET")
	r.HandleFunc("/api/settings", saveSettingsHandler).Methods("POST")
	r.HandleFunc("/ws/{id}", wsHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on http://localhost%s", addr)
	log.Printf("访问 http://localhost%s 使用A股智能分析系统", addr)
	
	// Open browser automatically
	go openBrowser(fmt.Sprintf("http://localhost%s", addr))
	
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>A股智能分析系统</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Microsoft YaHei', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        .header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
            animation: fadeInDown 0.6s ease;
            position: relative;
        }

        .header h1 {
            font-size: 3em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }

        .settings-btn {
            position: absolute;
            top: 0;
            right: 0;
            padding: 12px 24px;
            background: rgba(255,255,255,0.2);
            color: white;
            text-decoration: none;
            border-radius: 10px;
            font-size: 1em;
            font-weight: 600;
            transition: all 0.3s ease;
            display: inline-flex;
            align-items: center;
            gap: 8px;
        }

        .settings-btn:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.2);
        }

        .main-card {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            animation: fadeInUp 0.6s ease;
        }

        .form-group {
            margin-bottom: 25px;
        }

        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
            font-size: 1.1em;
        }

        .form-group input,
        .form-group select {
            width: 100%;
            padding: 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1em;
            transition: all 0.3s ease;
        }

        .form-group input:focus,
        .form-group select:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .button-group {
            display: flex;
            gap: 15px;
            margin-top: 30px;
        }

        .btn {
            flex: 1;
            padding: 15px 30px;
            font-size: 1.1em;
            font-weight: 600;
            border: none;
            border-radius: 10px;
            cursor: pointer;
            transition: all 0.3s ease;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
        }

        .btn-primary:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }

        .btn-secondary {
            background: #f0f0f0;
            color: #333;
        }

        .btn-secondary:hover {
            background: #e0e0e0;
        }

        .status-section {
            margin-top: 30px;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 10px;
            display: none;
        }

        .status-section.active {
            display: block;
            animation: fadeIn 0.3s ease;
        }

        .status-header {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }

        .status-icon {
            width: 40px;
            height: 40px;
            margin-right: 15px;
        }

        .status-icon.running {
            border: 3px solid #667eea;
            border-top-color: transparent;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }

        .status-icon.completed {
            color: #28a745;
            font-size: 40px;
        }

        .status-icon.failed {
            color: #dc3545;
            font-size: 40px;
        }

        .status-text {
            font-size: 1.2em;
            font-weight: 600;
        }

        .output-box {
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
            border-radius: 8px;
            max-height: 400px;
            overflow-y: auto;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 0.9em;
            line-height: 1.6;
            margin-top: 15px;
        }

        .output-line {
            margin-bottom: 5px;
        }

        .result-box {
            background: white;
            padding: 25px;
            border-radius: 10px;
            margin-top: 20px;
            border: 2px solid #e0e0e0;
        }

        .result-box h3 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.5em;
        }

        .result-content {
            white-space: pre-wrap;
            line-height: 1.8;
            color: #333;
        }

        .examples {
            margin-top: 20px;
            padding: 20px;
            background: #f0f7ff;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }

        .examples h3 {
            color: #667eea;
            margin-bottom: 15px;
        }

        .example-item {
            margin-bottom: 10px;
            padding: 10px;
            background: white;
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.2s ease;
        }

        .example-item:hover {
            transform: translateX(5px);
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }

        .example-item strong {
            color: #667eea;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        @keyframes fadeInDown {
            from {
                opacity: 0;
                transform: translateY(-20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @keyframes fadeInUp {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .progress-bar {
            width: 100%;
            height: 4px;
            background: #e0e0e0;
            border-radius: 2px;
            overflow: hidden;
            margin-top: 15px;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea, #764ba2);
            width: 0%;
            transition: width 0.3s ease;
            animation: progress 2s ease-in-out infinite;
        }

        @keyframes progress {
            0% { width: 0%; }
            50% { width: 70%; }
            100% { width: 100%; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <a href="/settings" class="settings-btn">⚙️ AI配置</a>
            <h1>🎯 A股智能分析系统</h1>
            <p>基于AI的专业股票分析平台</p>
        </div>

        <div class="main-card">
            <form id="analysisForm">
                <div class="form-group">
                    <label for="companyName">📊 公司名称</label>
                    <input type="text" id="companyName" name="company_name" placeholder="例如：贵州茅台" required>
                </div>

                <div class="form-group">
                    <label for="stockCode">🔢 股票代码</label>
                    <input type="text" id="stockCode" name="stock_code" placeholder="例如：600519.SH" required>
                </div>

                <div class="form-group">
                    <label for="market">🏛️ 市场类型</label>
                    <select id="market" name="market" required>
                        <option value="SH">上交所 (SH)</option>
                        <option value="SZ">深交所 (SZ)</option>
                        <option value="HK">港股 (HK)</option>
                    </select>
                </div>

                <div class="button-group">
                    <button type="submit" class="btn btn-primary" id="analyzeBtn">
                        🚀 开始分析
                    </button>
                    <button type="button" class="btn btn-secondary" id="resetBtn">
                        🔄 重置
                    </button>
                </div>
            </form>

            <div class="examples">
                <h3>📝 示例股票</h3>
                <div class="example-item" onclick="fillExample('贵州茅台', '600519.SH', 'SH')">
                    <strong>贵州茅台</strong> - 600519.SH (上交所)
                </div>
                <div class="example-item" onclick="fillExample('平安银行', '000001.SZ', 'SZ')">
                    <strong>平安银行</strong> - 000001.SZ (深交所)
                </div>
                <div class="example-item" onclick="fillExample('腾讯控股', '00700.HK', 'HK')">
                    <strong>腾讯控股</strong> - 00700.HK (港股)
                </div>
            </div>

            <div id="statusSection" class="status-section">
                <div class="status-header">
                    <div id="statusIcon" class="status-icon"></div>
                    <div id="statusText" class="status-text"></div>
                </div>
                <div class="progress-bar" id="progressBar" style="display:none;">
                    <div class="progress-fill"></div>
                </div>
                <div id="outputBox" class="output-box"></div>
                <div id="resultBox" class="result-box" style="display:none;">
                    <h3>📈 分析结果</h3>
                    <div id="resultContent" class="result-content"></div>
                </div>
            </div>
        </div>
    </div>

    <script>
        let currentSessionId = null;
        let ws = null;

        function fillExample(companyName, stockCode, market) {
            document.getElementById('companyName').value = companyName;
            document.getElementById('stockCode').value = stockCode;
            document.getElementById('market').value = market;
        }

        document.getElementById('resetBtn').addEventListener('click', function() {
            document.getElementById('analysisForm').reset();
            document.getElementById('statusSection').classList.remove('active');
            if (ws) {
                ws.close();
                ws = null;
            }
        });

        document.getElementById('analysisForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = {
                company_name: document.getElementById('companyName').value,
                stock_code: document.getElementById('stockCode').value,
                market: document.getElementById('market').value
            };

            // Show status section
            const statusSection = document.getElementById('statusSection');
            statusSection.classList.add('active');

            // Update status
            const statusIcon = document.getElementById('statusIcon');
            const statusText = document.getElementById('statusText');
            const outputBox = document.getElementById('outputBox');
            const resultBox = document.getElementById('resultBox');
            const progressBar = document.getElementById('progressBar');

            statusIcon.className = 'status-icon running';
            statusText.textContent = '正在启动分析...';
            outputBox.innerHTML = '';
            resultBox.style.display = 'none';
            progressBar.style.display = 'block';

            // Disable submit button
            const analyzeBtn = document.getElementById('analyzeBtn');
            analyzeBtn.disabled = true;
            analyzeBtn.textContent = '⏳ 分析中...';

            try {
                // Start analysis
                const response = await fetch('/api/analyze', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });

                const data = await response.json();
                currentSessionId = data.session_id;

                // Connect to WebSocket for real-time updates
                const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                ws = new WebSocket(protocol + '//' + window.location.host + '/ws/' + currentSessionId);

                ws.onmessage = function(event) {
                    const message = JSON.parse(event.data);
                    
                    if (message.type === 'output') {
                        const line = document.createElement('div');
                        line.className = 'output-line';
                        line.textContent = message.data;
                        outputBox.appendChild(line);
                        outputBox.scrollTop = outputBox.scrollHeight;
                    } else if (message.type === 'status') {
                        statusText.textContent = message.data;
                    } else if (message.type === 'completed') {
                        statusIcon.className = 'status-icon completed';
                        statusIcon.textContent = '✅';
                        statusText.textContent = '分析完成！';
                        progressBar.style.display = 'none';
                        
                        if (message.result) {
                            resultBox.style.display = 'block';
                            document.getElementById('resultContent').textContent = message.result;
                        }

                        analyzeBtn.disabled = false;
                        analyzeBtn.textContent = '🚀 开始分析';
                    } else if (message.type === 'error') {
                        statusIcon.className = 'status-icon failed';
                        statusIcon.textContent = '❌';
                        statusText.textContent = '分析失败：' + message.data;
                        progressBar.style.display = 'none';
                        
                        analyzeBtn.disabled = false;
                        analyzeBtn.textContent = '🚀 开始分析';
                    }
                };

                ws.onerror = function(error) {
                    console.error('WebSocket error:', error);
                    statusIcon.className = 'status-icon failed';
                    statusIcon.textContent = '❌';
                    statusText.textContent = '连接错误，请刷新页面重试';
                    progressBar.style.display = 'none';
                    analyzeBtn.disabled = false;
                    analyzeBtn.textContent = '🚀 开始分析';
                };

                ws.onclose = function() {
                    console.log('WebSocket connection closed');
                };

            } catch (error) {
                console.error('Error:', error);
                statusIcon.className = 'status-icon failed';
                statusIcon.textContent = '❌';
                statusText.textContent = '请求失败：' + error.message;
                progressBar.style.display = 'none';
                analyzeBtn.disabled = false;
                analyzeBtn.textContent = '🚀 开始分析';
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func settingsPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI 配置设置 - A股智能分析系统</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Microsoft YaHei', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 800px;
            margin: 0 auto;
        }

        .header {
            text-align: center;
            color: white;
            margin-bottom: 30px;
        }

        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .header p {
            font-size: 1.1em;
            opacity: 0.9;
        }

        .nav-back {
            display: inline-block;
            margin-bottom: 20px;
            padding: 10px 20px;
            background: rgba(255,255,255,0.2);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            transition: all 0.3s ease;
        }

        .nav-back:hover {
            background: rgba(255,255,255,0.3);
            transform: translateX(-5px);
        }

        .main-card {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
        }

        .form-group {
            margin-bottom: 25px;
        }

        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
            font-size: 1.1em;
        }

        .form-group .label-desc {
            font-size: 0.9em;
            color: #666;
            font-weight: normal;
            margin-top: 4px;
        }

        .form-group input,
        .form-group select {
            width: 100%;
            padding: 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1em;
            transition: all 0.3s ease;
        }

        .form-group input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .form-group input[type="password"] {
            font-family: monospace;
        }

        .button-group {
            display: flex;
            gap: 15px;
            margin-top: 30px;
        }

        .btn {
            flex: 1;
            padding: 15px 30px;
            font-size: 1.1em;
            font-weight: 600;
            border: none;
            border-radius: 10px;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #f0f0f0;
            color: #333;
        }

        .btn-secondary:hover {
            background: #e0e0e0;
        }

        .alert {
            padding: 15px 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            display: none;
        }

        .alert.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .alert.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .alert.show {
            display: block;
            animation: slideIn 0.3s ease;
        }

        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .info-box {
            background: #f0f7ff;
            padding: 20px;
            border-radius: 10px;
            border-left: 4px solid #667eea;
            margin-bottom: 25px;
        }

        .info-box h3 {
            color: #667eea;
            margin-bottom: 10px;
        }

        .info-box ul {
            margin-left: 20px;
            line-height: 1.8;
        }

        .loader {
            display: none;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin: 0 auto;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="nav-back">← 返回主页</a>
        
        <div class="header">
            <h1>⚙️ AI 配置设置</h1>
            <p>配置您的AI模型参数</p>
        </div>

        <div class="main-card">
            <div id="alertBox" class="alert"></div>

            <div class="info-box">
                <h3>📝 配置说明</h3>
                <ul>
                    <li><strong>API Key</strong>: 您的OpenAI API密钥</li>
                    <li><strong>Base URL</strong>: API服务地址（如使用代理或自建服务）</li>
                    <li><strong>Model Name</strong>: 使用的模型名称（如 gpt-4o, gpt-4, gpt-3.5-turbo）</li>
                    <li><strong>Temperature</strong>: 控制输出随机性（0-1，越高越随机）</li>
                    <li><strong>Max Tokens</strong>: 单次生成的最大token数</li>
                </ul>
            </div>

            <form id="settingsForm">
                <div class="form-group">
                    <label for="apiKey">
                        🔑 API Key
                        <div class="label-desc">您的OpenAI API密钥</div>
                    </label>
                    <input type="password" id="apiKey" name="api_key" placeholder="sk-..." required>
                </div>

                <div class="form-group">
                    <label for="baseUrl">
                        🌐 Base URL
                        <div class="label-desc">API服务地址</div>
                    </label>
                    <input type="text" id="baseUrl" name="base_url" placeholder="https://api.openai.com/v1" required>
                </div>

                <div class="form-group">
                    <label for="modelName">
                        🤖 Model Name
                        <div class="label-desc">使用的AI模型名称</div>
                    </label>
                    <input type="text" id="modelName" name="model_name" placeholder="gpt-4o" required>
                </div>

                <div class="form-group">
                    <label for="temperature">
                        🌡️ Temperature
                        <div class="label-desc">控制输出随机性（0.0 - 1.0）</div>
                    </label>
                    <input type="number" id="temperature" name="temperature" min="0" max="1" step="0.1" placeholder="0.8" required>
                </div>

                <div class="form-group">
                    <label for="maxTokens">
                        📊 Max Tokens
                        <div class="label-desc">单次生成的最大token数</div>
                    </label>
                    <input type="number" id="maxTokens" name="max_tokens" min="1000" max="128000" step="1000" placeholder="14000" required>
                </div>

                <div class="button-group">
                    <button type="submit" class="btn btn-primary" id="saveBtn">
                        💾 保存设置
                    </button>
                    <button type="button" class="btn btn-secondary" id="resetBtn">
                        🔄 重置为默认
                    </button>
                </div>
            </form>

            <div class="loader" id="loader"></div>
        </div>
    </div>

    <script>
        // Load settings on page load
        async function loadSettings() {
            try {
                const response = await fetch('/api/settings');
                const settings = await response.json();
                
                document.getElementById('apiKey').value = settings.api_key || '';
                document.getElementById('baseUrl').value = settings.base_url || 'https://api.openai.com/v1';
                document.getElementById('modelName').value = settings.model_name || 'gpt-4o';
                document.getElementById('temperature').value = settings.temperature || 0.8;
                document.getElementById('maxTokens').value = settings.max_tokens || 14000;
            } catch (error) {
                console.error('Error loading settings:', error);
            }
        }

        // Show alert message
        function showAlert(message, type) {
            const alertBox = document.getElementById('alertBox');
            alertBox.textContent = message;
            alertBox.className = 'alert ' + type + ' show';
            
            setTimeout(() => {
                alertBox.classList.remove('show');
            }, 5000);
        }

        // Save settings
        document.getElementById('settingsForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const saveBtn = document.getElementById('saveBtn');
            const loader = document.getElementById('loader');
            
            saveBtn.disabled = true;
            loader.style.display = 'block';
            
            const formData = {
                api_key: document.getElementById('apiKey').value,
                base_url: document.getElementById('baseUrl').value,
                model_name: document.getElementById('modelName').value,
                temperature: parseFloat(document.getElementById('temperature').value),
                max_tokens: parseInt(document.getElementById('maxTokens').value)
            };
            
            try {
                const response = await fetch('/api/settings', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    showAlert('✅ ' + result.message, 'success');
                } else {
                    showAlert('❌ 保存失败: ' + result.error, 'error');
                }
            } catch (error) {
                showAlert('❌ 保存失败: ' + error.message, 'error');
            } finally {
                saveBtn.disabled = false;
                loader.style.display = 'none';
            }
        });

        // Reset to default
        document.getElementById('resetBtn').addEventListener('click', function() {
            document.getElementById('apiKey').value = '';
            document.getElementById('baseUrl').value = 'https://api.openai.com/v1';
            document.getElementById('modelName').value = 'gpt-4o';
            document.getElementById('temperature').value = '0.8';
            document.getElementById('maxTokens').value = '14000';
            showAlert('ℹ️ 已重置为默认值，请点击"保存设置"以应用', 'success');
        });

        // Load settings when page loads
        loadSettings();
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new session
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
	session := &AnalysisSession{
		ID:          sessionID,
		Status:      "running",
		Output:      []string{},
		StartTime:   time.Now(),
		CompanyName: req.CompanyName,
		StockCode:   req.StockCode,
		Market:      req.Market,
		subscribers: []*websocket.Conn{},
	}

	sessionMu.Lock()
	sessions[sessionID] = session
	sessionMu.Unlock()

	// Start analysis in background
	go runAnalysis(session)

	// Return session ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": sessionID,
	})
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	sessionMu.RLock()
	session, exists := sessions[sessionID]
	sessionMu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	session.mu.RLock()
	defer session.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	sessionMu.RLock()
	session, exists := sessions[sessionID]
	sessionMu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	session.mu.Lock()
	session.subscribers = append(session.subscribers, conn)
	session.mu.Unlock()

	// Send existing output
	session.mu.RLock()
	for _, line := range session.Output {
		msg := map[string]string{
			"type": "output",
			"data": line,
		}
		conn.WriteJSON(msg)
	}
	session.mu.RUnlock()

	// Keep connection alive
	defer func() {
		session.mu.Lock()
		// Remove connection from subscribers
		for i, sub := range session.subscribers {
			if sub == conn {
				session.subscribers = append(session.subscribers[:i], session.subscribers[i+1:]...)
				break
			}
		}
		session.mu.Unlock()
		conn.Close()
	}()

	// Read messages from client (keep-alive)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func runAnalysis(session *AnalysisSession) {
	defer func() {
		session.mu.Lock()
		session.EndTime = time.Now()
		session.mu.Unlock()
	}()

	session.broadcastStatus("正在启动Python分析引擎...")

	// Try to find bundled Python engine first
	enginePath := findBundledEngine()
	var cmd *exec.Cmd
	
	if enginePath != "" {
		// Use bundled engine
		session.broadcastStatus("使用内置Python引擎...")
		cmd = exec.Command(
			enginePath,
			"--company", session.CompanyName,
			"--code", session.StockCode,
			"--market", session.Market,
		)
	} else {
		// Fallback to traditional Python method
		session.broadcastStatus("使用系统Python环境...")

		pythonCmd := findPython()
		if pythonCmd == "" {
			session.broadcastError("未找到Python环境，请确保已安装Python 3.12+")
			return
		}

		cliEntryPath := findCliEntryPy()
		if cliEntryPath == "" {
			session.broadcastError("未找到Python分析脚本 (cli_entry.py)")
			return
		}

		// Use the cli_entry.py script directly
		cmd = exec.Command(
			pythonCmd,
			"-u", // Add -u flag to disable output buffering
			cliEntryPath,
			"--company", session.CompanyName,
			"--code", session.StockCode,
			"--market", session.Market,
		)
		// Set the working directory to the script's directory
		cmd.Dir = filepath.Dir(cliEntryPath)
	}

	// Set up pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		session.broadcastError(fmt.Sprintf("创建stdout管道失败: %v", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		session.broadcastError(fmt.Sprintf("创建stderr管道失败: %v", err))
		return
	}

	if err := cmd.Start(); err != nil {
		errMsg := fmt.Sprintf("启动Python进程失败: %v", err)
		log.Printf(errMsg)
		session.broadcastError(errMsg)
		return
	}

	log.Printf("Python process started successfully, PID: %d", cmd.Process.Pid)
	session.broadcastStatus("Python分析引擎已启动，正在加载分析模块...")

	// Read output in goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	var resultBuilder string
	var outputMu sync.Mutex

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			// Process structured messages
			if strings.HasPrefix(line, "[STATUS]") {
				statusMsg := strings.TrimSpace(strings.TrimPrefix(line, "[STATUS]"))
				session.broadcastStatus(statusMsg)
			} else if strings.HasPrefix(line, "[PROGRESS]") {
				// Progress messages can be handled differently if needed
				// For now, just treat them as regular output
				session.broadcastOutput(line)
			} else {
				session.broadcastOutput(line)
			}

			outputMu.Lock()
			resultBuilder += line + "\n"
			outputMu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			// 将stderr也输出到界面，但加上标记
			session.broadcastOutput("⚠️ " + line)
			log.Printf("Python stderr: %s", line)
		}
	}()

	wg.Wait()
	err = cmd.Wait()

	session.mu.Lock()
	session.Result = resultBuilder
	session.mu.Unlock()

	if err != nil {
		errMsg := fmt.Sprintf("分析过程出错: %v", err)
		log.Printf("Analysis failed: %v", err)
		session.broadcastError(errMsg)
	} else {
		log.Printf("Analysis completed successfully")
		session.broadcastComplete(resultBuilder)
	}
}

func (s *AnalysisSession) broadcastOutput(line string) {
	s.mu.Lock()
	s.Output = append(s.Output, line)
	subscribers := make([]*websocket.Conn, len(s.subscribers))
	copy(subscribers, s.subscribers)
	s.mu.Unlock()

	msg := map[string]string{
		"type": "output",
		"data": line,
	}

	for _, conn := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing to websocket: %v", err)
		}
	}
}

func (s *AnalysisSession) broadcastStatus(status string) {
	s.mu.RLock()
	subscribers := make([]*websocket.Conn, len(s.subscribers))
	copy(subscribers, s.subscribers)
	s.mu.RUnlock()

	msg := map[string]string{
		"type": "status",
		"data": status,
	}

	for _, conn := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing to websocket: %v", err)
		}
	}
}

func (s *AnalysisSession) broadcastComplete(result string) {
	s.mu.Lock()
	s.Status = "completed"
	s.mu.Unlock()

	s.mu.RLock()
	subscribers := make([]*websocket.Conn, len(s.subscribers))
	copy(subscribers, s.subscribers)
	s.mu.RUnlock()

	msg := map[string]interface{}{
		"type":   "completed",
		"result": result,
	}

	for _, conn := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing to websocket: %v", err)
		}
	}
}

func (s *AnalysisSession) broadcastError(errMsg string) {
	s.mu.Lock()
	s.Status = "failed"
	s.mu.Unlock()

	s.mu.RLock()
	subscribers := make([]*websocket.Conn, len(s.subscribers))
	copy(subscribers, s.subscribers)
	s.mu.RUnlock()

	msg := map[string]string{
		"type": "error",
		"data": errMsg,
	}

	for _, conn := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing to websocket: %v", err)
		}
	}
}

func findBundledEngine() string {
	// Get the directory where the executable is located
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	exeDir := filepath.Dir(exePath)
	
	// Check for bundled engine in several locations
	possiblePaths := []string{
		// Next to the executable
		filepath.Join(exeDir, "python_bundle", "stock_analysis_engine", "stock_analysis_engine.exe"),
		filepath.Join(exeDir, "python_bundle", "stock_analysis_engine", "stock_analysis_engine"),
		// In a subdirectory
		filepath.Join(exeDir, "stock_analysis_engine", "stock_analysis_engine.exe"),
		filepath.Join(exeDir, "stock_analysis_engine", "stock_analysis_engine"),
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}

func findPython() string {
	// Try common Python commands
	commands := []string{"python3", "python", "py"}
	
	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			// Verify it's Python 3
			output, err := exec.Command(path, "--version").Output()
			if err == nil && len(output) > 0 {
				return path
			}
		}
	}
	
	return ""
}

func findCliEntryPy() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Try to find cli_entry.py in stock_analysis_a_stock directory
	possiblePaths := []string{
		filepath.Join(cwd, "..", "stock_analysis_a_stock", "src", "a_stock_analysis", "cli_entry.py"),
		filepath.Join(cwd, "..", "..", "stock_analysis_a_stock", "src", "a_stock_analysis", "cli_entry.py"),
		filepath.Join(cwd, "stock_analysis_a_stock", "src", "a_stock_analysis", "cli_entry.py"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func openBrowser(url string) {
	time.Sleep(1 * time.Second) // Wait for server to start
	
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	
	if err != nil {
		log.Printf("无法自动打开浏览器: %v", err)
		log.Printf("请手动访问: %s", url)
	}
}
