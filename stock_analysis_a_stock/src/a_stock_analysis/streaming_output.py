#!/usr/bin/env python3
"""
实时输出捕获模块
用于确保所有输出都能实时传递到UI
"""
import sys
import io
from contextlib import redirect_stdout, redirect_stderr


class FlushingTextIOWrapper(io.TextIOWrapper):
    """自动刷新的文本IO包装器"""
    
    def write(self, s):
        result = super().write(s)
        self.flush()
        return result


def setup_realtime_output():
    """
    配置实时输出
    确保所有print语句和日志输出都能立即刷新
    """
    # 重新配置stdout和stderr为行缓冲模式
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(line_buffering=True)
    if hasattr(sys.stderr, 'reconfigure'):
        sys.stderr.reconfigure(line_buffering=True)
    
    # 设置环境变量禁用Python缓冲
    import os
    os.environ['PYTHONUNBUFFERED'] = '1'


def flush_all():
    """刷新所有输出缓冲区"""
    sys.stdout.flush()
    sys.stderr.flush()
