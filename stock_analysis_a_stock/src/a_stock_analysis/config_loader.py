"""
Configuration loader for AI settings.
Loads settings from .env file or from user's settings file (created by UI).
"""
import os
import json
from pathlib import Path
from typing import Dict, Any


def get_settings_file_path() -> Path:
    """Get the path to the settings file."""
    home_dir = Path.home()
    settings_dir = home_dir / ".a_stock_analysis"
    settings_dir.mkdir(exist_ok=True)
    return settings_dir / "ai_settings.json"


def load_ai_config() -> Dict[str, Any]:
    """
    Load AI configuration.
    Priority: settings file > environment variables > defaults
    """
    config = {
        "model_name": "gpt-4o",
        "api_key": "",
        "base_url": "https://api.openai.com/v1",
        "temperature": 0.8,
        "max_tokens": 14000,
    }
    
    # First, try to load from environment variables
    if os.getenv("OPENAI_API_KEY"):
        config["api_key"] = os.getenv("OPENAI_API_KEY")
    if os.getenv("OPENAI_BASE_URL"):
        config["base_url"] = os.getenv("OPENAI_BASE_URL")
    if os.getenv("OPENAI_MODEL_NAME"):
        config["model_name"] = os.getenv("OPENAI_MODEL_NAME")
    if os.getenv("TEMPERATURE"):
        try:
            config["temperature"] = float(os.getenv("TEMPERATURE"))
        except ValueError:
            pass
    if os.getenv("MAX_TOKENS"):
        try:
            config["max_tokens"] = int(os.getenv("MAX_TOKENS"))
        except ValueError:
            pass
    
    # Then, try to load from settings file (overrides env vars)
    settings_path = get_settings_file_path()
    if settings_path.exists():
        try:
            with open(settings_path, 'r', encoding='utf-8') as f:
                settings = json.load(f)
                if settings.get("api_key"):
                    config["api_key"] = settings["api_key"]
                if settings.get("base_url"):
                    config["base_url"] = settings["base_url"]
                if settings.get("model_name"):
                    config["model_name"] = settings["model_name"]
                if settings.get("temperature") is not None:
                    config["temperature"] = settings["temperature"]
                if settings.get("max_tokens") is not None:
                    config["max_tokens"] = settings["max_tokens"]
        except (json.JSONDecodeError, IOError) as e:
            print(f"警告: 无法加载设置文件: {e}")
    
    return config
