#!/usr/bin/env python3
"""
闲鱼 Deeplink 提取 CLI 脚本
供 Go 服务通过 exec.Command 调用

用法:
  # 单链模式
  python3 xianyu_deeplink.py --mode single --url "https://m.tb.cn/h.xxx" --platform android [--driver /path/to/chromedriver]

  # 批量模式
  python3 xianyu_deeplink.py --mode batch --input /tmp/input.txt --output /tmp/output.xlsx --platform android [--driver /path/to/chromedriver]

输出: 标准 JSON（stdout）
"""

import sys
import os
import json
import argparse
import concurrent.futures
import queue

# 将 python-automation-project/src 加入 Python 路径
_SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
_SRC_DIR = os.path.join(_SCRIPT_DIR, '..', '..', 'python-automation-project', 'src')
sys.path.insert(0, os.path.normpath(_SRC_DIR))

from extract_xianyu_deeplink import get_xianyu_deeplink, CHROME_DRIVER_PATH as _DEFAULT_DRIVER


def _get_driver_path(arg_driver):
    if arg_driver:
        return arg_driver
    return _DEFAULT_DRIVER


def single_mode(url, platform, driver_path):
    """单链模式：提取一条闲鱼链接的 deeplink"""
    import extract_xianyu_deeplink as mod
    mod.CHROME_DRIVER_PATH = driver_path

    try:
        deeplink = get_xianyu_deeplink(url, platform=platform)
        if deeplink:
            print(json.dumps({
                "deeplink": deeplink,
                "error": None
            }, ensure_ascii=False))
        else:
            print(json.dumps({
                "deeplink": "",
                "error": "未能提取 deeplink，请检查链接是否有效"
            }, ensure_ascii=False))
    except Exception as e:
        print(json.dumps({
            "deeplink": "",
            "error": str(e)
        }, ensure_ascii=False))


def batch_mode(input_file, output_file, platform, driver_path):
    """批量模式：从 txt 文件批量提取闲鱼 deeplink，输出 Excel"""
    import extract_xianyu_deeplink as mod
    mod.CHROME_DRIVER_PATH = driver_path

    from selenium import webdriver
    from selenium.webdriver.chrome.service import Service
    from selenium.webdriver.chrome.options import Options
    import openpyxl

    # 读取输入文件
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            short_urls = [line.strip() for line in f if line.strip()]
    except Exception as e:
        print(json.dumps({"success": 0, "failed": 0, "output": "", "error": f"读取输入文件失败: {e}"}))
        return

    if not short_urls:
        print(json.dumps({"success": 0, "failed": 0, "output": "", "error": "输入文件为空"}))
        return

    # 创建驱动池
    driver_num = min(3, len(short_urls))
    driver_pool = queue.Queue()

    def create_driver():
        chrome_options = Options()
        if platform.lower() == "android":
            mobile_emulation = {
                "deviceMetrics": {"width": 412, "height": 915, "pixelRatio": 3.5},
                "userAgent": "Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36"
            }
        else:
            mobile_emulation = {
                "deviceMetrics": {"width": 375, "height": 812, "pixelRatio": 3.0},
                "userAgent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
            }
        chrome_options.add_experimental_option("mobileEmulation", mobile_emulation)
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument('log-level=3')
        try:
            service = Service(driver_path)
            return webdriver.Chrome(service=service, options=chrome_options)
        except Exception:
            return webdriver.Chrome(options=chrome_options)

    drivers = []
    for _ in range(driver_num):
        try:
            d = create_driver()
            drivers.append(d)
            driver_pool.put(d)
        except Exception:
            pass

    results = []
    success_count = 0
    failed_count = 0

    def process_url(short_url):
        nonlocal success_count, failed_count
        driver = driver_pool.get()
        try:
            deeplink = get_xianyu_deeplink(short_url, driver=driver, platform=platform)
            if deeplink:
                success_count += 1
                return {"原始链接": short_url, "Deeplink": deeplink, "状态": "成功"}
            else:
                failed_count += 1
                return {"原始链接": short_url, "Deeplink": "", "状态": "失败"}
        except Exception as e:
            failed_count += 1
            return {"原始链接": short_url, "Deeplink": "", "状态": f"失败: {e}"}
        finally:
            driver_pool.put(driver)

    with concurrent.futures.ThreadPoolExecutor(max_workers=min(2, driver_num)) as executor:
        futures = {executor.submit(process_url, url): url for url in short_urls}
        for future in concurrent.futures.as_completed(futures):
            results.append(future.result())

    for d in drivers:
        try:
            d.quit()
        except Exception:
            pass

    # 写入 Excel
    try:
        wb = openpyxl.Workbook()
        ws = wb.active
        ws.title = "闲鱼Deeplink结果"
        headers = ["原始链接", "Deeplink", "状态"]
        ws.append(headers)
        for row in results:
            ws.append([row.get(h, "") for h in headers])
        wb.save(output_file)
    except Exception as e:
        print(json.dumps({"success": success_count, "failed": failed_count, "output": "", "error": f"写入 Excel 失败: {e}"}))
        return

    print(json.dumps({
        "success": success_count,
        "failed": failed_count,
        "output": output_file,
        "error": None
    }, ensure_ascii=False))


def main():
    parser = argparse.ArgumentParser(description="闲鱼 Deeplink 提取工具")
    parser.add_argument('--mode', choices=['single', 'batch'], required=True, help="运行模式")
    parser.add_argument('--url', help="单链模式：闲鱼短链接")
    parser.add_argument('--platform', default='android', choices=['android', 'ios'], help="平台 (默认 android)")
    parser.add_argument('--input', help="批量模式：输入 txt 文件路径")
    parser.add_argument('--output', help="批量模式：输出 Excel 文件路径")
    parser.add_argument('--driver', help="ChromeDriver 可执行文件路径（可选）")
    args = parser.parse_args()

    driver_path = _get_driver_path(args.driver)

    if args.mode == 'single':
        if not args.url:
            print(json.dumps({"deeplink": "", "error": "--url 参数不能为空"}))
            sys.exit(1)
        single_mode(args.url, args.platform, driver_path)
    elif args.mode == 'batch':
        if not args.input or not args.output:
            print(json.dumps({"success": 0, "failed": 0, "output": "", "error": "--input 和 --output 参数不能为空"}))
            sys.exit(1)
        batch_mode(args.input, args.output, args.platform, driver_path)


if __name__ == '__main__':
    main()
