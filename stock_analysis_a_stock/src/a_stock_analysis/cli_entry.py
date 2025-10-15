#!/usr/bin/env python3
"""
CLI entry point for A股分析系统
用于PyInstaller打包和命令行调用
"""
import sys
import os
import argparse

# 首先设置实时输出
from streaming_output import setup_realtime_output, flush_all
setup_realtime_output()

from crew import AStockAnalysisCrew

def flush_print(*args, **kwargs):
    """立即刷新输出的print函数"""
    print(*args, **kwargs)
    flush_all()

def main():
    """主入口函数"""
    
    parser = argparse.ArgumentParser(description='A股智能分析系统')
    parser.add_argument('--company', type=str, required=True, help='公司名称')
    parser.add_argument('--code', type=str, required=True, help='股票代码')
    parser.add_argument('--market', type=str, required=True, help='市场类型 (SH/SZ/HK)')
    
    args = parser.parse_args()
    
    inputs = {
        'company_name': args.company,
        'stock_code': args.code,
        'market': args.market
    }
    
    flush_print("## 欢迎使用A股智能分析系统")
    flush_print('-------------------------------')
    flush_print(f"正在分析: {inputs['company_name']} ({inputs['stock_code']})")
    flush_print('-------------------------------')
    flush_print("[STATUS] 正在初始化AI分析团队...")
    
    try:
        flush_print("[STATUS] 启动分析流程...")
        result = AStockAnalysisCrew().crew().kickoff(inputs=inputs)
        
        flush_print("\n\n########################")
        flush_print("## 分析报告")
        flush_print("########################\n")
        flush_print(str(result))
        flush_print("\n[STATUS] 分析完成！")
        
        return 0
    except Exception as e:
        error_msg = f"\n错误: {e}"
        print(error_msg, file=sys.stderr)
        sys.stderr.flush()
        return 1

if __name__ == "__main__":
    sys.exit(main())
