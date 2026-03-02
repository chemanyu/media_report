#!/usr/bin/env python3
"""
淘宝客活动短链批量获取 CLI 脚本
供 Go 服务通过 exec.Command 调用

用法:
  python3 activity_batch.py \
    --input /tmp/pids.txt \
    --output /tmp/output.xlsx \
    --material_id 20150318020010092 \
    --app_key 35238422 \
    --app_secret 3e2c5266e7a3689ac909659a203ce301 \
    [--driver /path/to/chromedriver]

输入文件格式（txt，每行一个 sub_pid）:
  mm_874030133_3340450211_116193450327
  mm_874030133_3340450211_116193450328

输出: 标准 JSON（stdout）
"""

import sys
import os
import json
import argparse
import concurrent.futures

# 将 python-automation-project/src 加入 Python 路径
_SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
_SRC_DIR = os.path.join(_SCRIPT_DIR, '..', '..', 'python-automation-project', 'src')
sys.path.insert(0, os.path.normpath(_SRC_DIR))

from taobao_api import TaobaoAPI
from extract_taobao_deeplink import get_taobao_deeplink, CHROME_DRIVER_PATH as _DEFAULT_DRIVER


def _get_driver_path(arg_driver):
    if arg_driver:
        return arg_driver
    return _DEFAULT_DRIVER


def process_sub_pid(sub_pid, api, driver_path, material_id):
    """处理单个 sub_pid：获取活动信息 → 转短链 → 提取 deeplink"""
    import extract_taobao_deeplink as mod
    mod.CHROME_DRIVER_PATH = driver_path

    # 提取 adzone_id（sub_pid 最后一个下划线分隔的部分）
    parts = sub_pid.strip().split('_')
    if len(parts) < 1:
        return {
            "推广位": sub_pid,
            "会场名称": "",
            "推广长链": "",
            "推广短链": "",
            "Deeplink": "",
            "h5Dp": "",
            "状态": "失败: sub_pid 格式错误"
        }

    adzone_id = parts[-1]

    # 1. 获取活动信息（推广长链）
    activity_info = api.get_activity_info(
        activity_material_id=material_id,
        adzone_id=adzone_id,
        sub_pid=sub_pid.strip()
    )

    if not activity_info:
        return {
            "推广位": sub_pid,
            "会场名称": "",
            "推广长链": "",
            "推广短链": "",
            "Deeplink": "",
            "h5Dp": "",
            "状态": "失败: 获取活动信息失败"
        }

    page_name = activity_info.get('page_name', '')
    click_url = activity_info.get('click_url', '')

    if not click_url:
        return {
            "推广位": sub_pid,
            "会场名称": page_name,
            "推广长链": "",
            "推广短链": "",
            "Deeplink": "",
            "h5Dp": "",
            "状态": "失败: 推广长链为空"
        }

    # 2. 将推广长链转为短链
    short_results = api.convert_to_short_url(click_url)
    short_url = ""
    if short_results and len(short_results) > 0:
        short_url = short_results[0].get('content', '')

    if not short_url:
        return {
            "推广位": sub_pid,
            "会场名称": page_name,
            "推广长链": click_url,
            "推广短链": "",
            "Deeplink": "",
            "h5Dp": "",
            "状态": "失败: 短链转换失败"
        }

    # 3. 提取 deeplink
    deeplink, h5_dp = get_taobao_deeplink(short_url, platform="ios")
    if deeplink:
        return {
            "推广位": sub_pid,
            "会场名称": page_name,
            "推广长链": click_url,
            "推广短链": short_url,
            "Deeplink": deeplink,
            "h5Dp": h5_dp or "",
            "状态": "成功"
        }
    else:
        return {
            "推广位": sub_pid,
            "会场名称": page_name,
            "推广长链": click_url,
            "推广短链": short_url,
            "Deeplink": "",
            "h5Dp": "",
            "状态": "失败: Deeplink 提取失败"
        }


def main():
    parser = argparse.ArgumentParser(description="淘宝客活动短链批量获取工具")
    parser.add_argument('--input', required=True, help="输入 txt 文件路径（每行一个 sub_pid）")
    parser.add_argument('--output', required=True, help="输出 Excel 文件路径")
    parser.add_argument('--material_id', default='20150318020010092', help="活动素材ID")
    parser.add_argument('--app_key', required=True, help="淘宝客 AppKey")
    parser.add_argument('--app_secret', required=True, help="淘宝客 AppSecret")
    parser.add_argument('--driver', help="ChromeDriver 路径（可选）")
    args = parser.parse_args()

    driver_path = _get_driver_path(args.driver)

    # 读取 sub_pid 列表
    try:
        with open(args.input, 'r', encoding='utf-8') as f:
            sub_pids = [line.strip() for line in f if line.strip()]
    except Exception as e:
        print(json.dumps({"success": 0, "failed": 0, "output": "", "error": f"读取输入文件失败: {e}"}))
        sys.exit(1)

    if not sub_pids:
        print(json.dumps({"success": 0, "failed": 0, "output": "", "error": "输入文件为空"}))
        sys.exit(1)

    api = TaobaoAPI(args.app_key, args.app_secret)

    results = []
    success_count = 0
    failed_count = 0

    with concurrent.futures.ThreadPoolExecutor(max_workers=2) as executor:
        futures = {
            executor.submit(process_sub_pid, pid, api, driver_path, args.material_id): pid
            for pid in sub_pids
        }
        for future in concurrent.futures.as_completed(futures):
            result = future.result()
            results.append(result)
            if result.get("状态") == "成功":
                success_count += 1
            else:
                failed_count += 1

    # 写入 Excel
    try:
        import openpyxl
        wb = openpyxl.Workbook()
        ws = wb.active
        ws.title = "活动短链结果"
        headers = ["推广位", "会场名称", "推广长链", "推广短链", "Deeplink", "h5Dp", "状态"]
        ws.append(headers)
        for row in results:
            ws.append([row.get(h, "") for h in headers])
        wb.save(args.output)
    except Exception as e:
        print(json.dumps({"success": success_count, "failed": failed_count, "output": "", "error": f"写入 Excel 失败: {e}"}))
        sys.exit(1)

    print(json.dumps({
        "success": success_count,
        "failed": failed_count,
        "output": args.output,
        "error": None
    }, ensure_ascii=False))


if __name__ == '__main__':
    main()
