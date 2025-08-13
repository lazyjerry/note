#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Git 自動添加推送到遠端倉庫工具
"""

import subprocess
import sys
import os
from typing import Optional, Tuple


def run_command(command: list, capture_output: bool = True) -> Tuple[int, str, str]:
    """
    執行命令並返回結果
    
    Args:
        command: 要執行的命令列表
        capture_output: 是否捕獲輸出
    
    Returns:
        (return_code, stdout, stderr)
    """
    try:
        result = subprocess.run(
            command,
            capture_output=capture_output,
            text=True,
            encoding='utf-8'
        )
        return result.returncode, result.stdout, result.stderr
    except Exception as e:
        print(f"執行命令時發生錯誤: {e}")
        return 1, "", str(e)


def check_git_repository() -> bool:
    """檢查當前目錄是否為 Git 倉庫"""
    return_code, _, _ = run_command(['git', 'rev-parse', '--git-dir'])
    return return_code == 0


def get_git_status() -> str:
    """獲取 Git 狀態"""
    return_code, stdout, stderr = run_command(['git', 'status', '--porcelain'])
    if return_code != 0:
        print(f"獲取 Git 狀態失敗: {stderr}")
        return ""
    return stdout


def add_all_files() -> bool:
    """添加所有變更的檔案"""
    print("正在添加所有變更的檔案...")
    return_code, stdout, stderr = run_command(['git', 'add', '.'])
    if return_code != 0:
        print(f"添加檔案失敗: {stderr}")
        return False
    print("檔案添加成功！")
    return True


def get_commit_message() -> Optional[str]:
    """獲取用戶輸入的 commit message"""
    print("\n" + "="*50)
    print("請輸入 commit message:")
    print("="*50)
    
    message = input().strip()
    
    if not message:
        print("錯誤: commit message 不能為空！")
        return None
    
    return message


def confirm_commit(message: str) -> bool:
    """確認是否要提交"""
    print("\n" + "="*50)
    print("確認提交資訊:")
    print(f"Commit Message: {message}")
    print("="*50)
    
    while True:
        confirm = input("是否確認提交？(y/n): ").strip().lower()
        if confirm in ['y', 'yes', '是', '確認']:
            return True
        elif confirm in ['n', 'no', '否', '取消']:
            return False
        else:
            print("請輸入 y 或 n")


def commit_changes(message: str) -> bool:
    """提交變更"""
    print("正在提交變更...")
    return_code, stdout, stderr = run_command(['git', 'commit', '-m', message])
    if return_code != 0:
        print(f"提交失敗: {stderr}")
        return False
    print("提交成功！")
    return True


def push_to_remote() -> bool:
    """推送到遠端倉庫"""
    print("正在推送到遠端倉庫...")
    
    # 獲取當前分支名稱
    return_code, branch, stderr = run_command(['git', 'branch', '--show-current'])
    if return_code != 0:
        print(f"獲取分支名稱失敗: {stderr}")
        return False
    
    branch = branch.strip()
    
    # 推送到遠端
    return_code, stdout, stderr = run_command(['git', 'push', 'origin', branch])
    if return_code != 0:
        print(f"推送失敗: {stderr}")
        return False
    
    print(f"成功推送到遠端分支: {branch}")
    return True


def main():
    """主函數"""
    print("Git 自動添加推送到遠端倉庫工具")
    print("="*50)
    
    # 檢查是否為 Git 倉庫
    if not check_git_repository():
        print("錯誤: 當前目錄不是 Git 倉庫！")
        print("請在 Git 倉庫目錄中執行此腳本。")
        sys.exit(1)
    
    # 檢查是否有變更
    status = get_git_status()
    if not status:
        print("沒有需要提交的變更。")
        sys.exit(0)
    
    print("檢測到以下變更:")
    print(status)
    
    # 添加所有檔案
    if not add_all_files():
        sys.exit(1)
    
    # 獲取 commit message
    message = get_commit_message()
    if not message:
        sys.exit(1)
    
    # 確認提交
    if not confirm_commit(message):
        print("已取消提交。")
        sys.exit(0)
    
    # 提交變更
    if not commit_changes(message):
        sys.exit(1)
    
    # 推送到遠端
    if not push_to_remote():
        sys.exit(1)
    
    print("\n" + "="*50)
    print("所有操作完成！")
    print("="*50)


if __name__ == "__main__":
    main()
