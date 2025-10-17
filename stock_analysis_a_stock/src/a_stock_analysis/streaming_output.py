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
    # 设置环境变量，强制 UTF-8 & 关闭缓冲，避免 Windows GBK 导致的编码报错与中文乱码
    import os
    os.environ['PYTHONUNBUFFERED'] = '1'
    # 优先通过环境变量强制 Python 使用 UTF-8（适用于 PyInstaller/子进程场景）
    os.environ.setdefault('PYTHONIOENCODING', 'utf-8')
    os.environ.setdefault('PYTHONUTF8', '1')
    # 在部分系统上，设置区域语言也有助于统一为 UTF-8
    os.environ.setdefault('LC_ALL', 'C.UTF-8')
    os.environ.setdefault('LANG', 'C.UTF-8')

    # 重新配置 stdout/stderr：行缓冲 + UTF-8 编码，并在出现无法编码字符时进行替换
    # 这样即使上游库输出 emoji 等字符，也不会因编码异常而崩溃
    if hasattr(sys.stdout, 'reconfigure'):
        try:
            sys.stdout.reconfigure(encoding='utf-8', errors='replace', line_buffering=True)
        except Exception:
            # 某些环境下 reconfigure 可能不可用或受限，忽略即可
            pass
    if hasattr(sys.stderr, 'reconfigure'):
        try:
            sys.stderr.reconfigure(encoding='utf-8', errors='replace', line_buffering=True)
        except Exception:
            pass

    # 兜底方案：直接以 UTF-8 包装底层 buffer，确保编码为 UTF-8（适用于 PyInstaller/管道场景）
    try:
        if hasattr(sys.stdout, 'buffer'):
            sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace', line_buffering=True)
    except Exception:
        pass
    try:
        if hasattr(sys.stderr, 'buffer'):
            sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8', errors='replace', line_buffering=True)
    except Exception:
        pass


def flush_all():
    """刷新所有输出缓冲区"""
    sys.stdout.flush()
    sys.stderr.flush()
