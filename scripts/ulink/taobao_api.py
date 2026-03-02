"""
淘宝客 API 调用模块
用于调用淘宝客活动素材信息查询接口
可以在其他模块中导入使用
"""
import hashlib
import time
import requests


class TaobaoAPI:
    """淘宝客 API 调用类"""

    def __init__(self, app_key, app_secret):
        """
        初始化淘宝客 API

        Args:
            app_key: TOP分配给应用的AppKey
            app_secret: TOP分配给应用的AppSecret（用于签名）
        """
        self.app_key = app_key
        self.app_secret = app_secret
        self.api_url = "http://gw.api.taobao.com/router/rest"

    def generate_sign(self, params):
        """
        生成 API 签名

        签名步骤：
        1. 按照ASCII顺序排序参数名与参数值
        2. 拼接参数名与参数值
        3. 使用md5加密
        4. 转换为大写

        Args:
            params: 请求参数字典

        Returns:
            签名字符串
        """
        # 1. 按ASCII顺序排序参数（去掉sign参数）
        sorted_params = sorted([(k, v) for k, v in params.items() if k != 'sign'])

        # 2. 拼接参数名与参数值
        sign_string = self.app_secret
        for key, value in sorted_params:
            sign_string += str(key) + str(value)
        sign_string += self.app_secret

        # 3. 使用md5加密
        md5_hash = hashlib.md5()
        md5_hash.update(sign_string.encode('utf-8'))

        # 4. 转换为大写
        sign = md5_hash.hexdigest().upper()

        return sign

    def get_activity_info(self, activity_material_id, adzone_id, sub_pid=None, relation_id=None, union_id=None):
        """
        获取淘宝客活动素材信息

        Args:
            activity_material_id: 官方活动会场ID，从淘宝客后台"活动推广"中获取
            adzone_id: mm_xxx_xxx_xxx的第三位
            sub_pid: mm_xxx_xxx_xxx 仅三方分成场景使用
            relation_id: 渠道关系id
            union_id: 自定义输入串，英文和数字组成，长度不能大于12个字符

        Returns:
            返回结果字典，包含：
            {
                'wx_qrcode_url': '微信群二维码地址',
                'click_url': '淘客推广长链',
                'short_click_url': '淘客推广短链',
                'terminal_type': '投放平台',
                'material_oss_url': '物料素材下载地址',
                'page_name': '会场名称',
                'page_start_time': '活动开始时间',
                'page_end_time': '活动结束时间'
            }
            失败返回 None
        """
        # 构建公共参数
        params = {
            'method': 'taobao.tbk.activity.info.get',
            'app_key': self.app_key,
            'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
            'format': 'json',
            'v': '2.0',
            'sign_method': 'md5',
            'partner_id': 'top-apitools',
        }

        # 添加业务参数
        params['activity_material_id'] = str(activity_material_id)
        params['adzone_id'] = str(adzone_id)

        if sub_pid:
            params['sub_pid'] = str(sub_pid)
        if relation_id:
            params['relation_id'] = str(relation_id)
        if union_id:
            params['union_id'] = str(union_id)

        # 生成签名
        params['sign'] = self.generate_sign(params)

        # 发起请求
        try:
            response = requests.get(self.api_url, params=params, timeout=10)
            response.raise_for_status()

            result = response.json()

            # 检查是否有错误
            if 'error_response' in result:
                error = result['error_response']
                print(f"淘宝客 API 错误: {error.get('msg', '未知错误')}")
                print(f"错误代码: {error.get('code', 'N/A')}")
                print(f"子错误代码: {error.get('sub_code', 'N/A')}")
                return None

            # 返回成功结果
            if 'tbk_activity_info_get_response' in result:
                data = result['tbk_activity_info_get_response'].get('data', {})
                return data

            return None

        except requests.exceptions.RequestException as e:
            print(f"淘宝客 API 请求失败: {e}")
            return None
        except Exception as e:
            print(f"淘宝客 API 处理异常: {e}")
            return None

    def convert_to_short_url(self, urls):
        """
        将淘宝客推广长链转换为短链

        Args:
            urls: 长链接列表，每个元素是一个字典 {"url": "长链接"}
                  或者直接传入字符串列表

        Returns:
            返回结果字典列表，每个元素包含：
            {
                'short_url': '短链接',
                'err_msg': '错误信息（如果有）'
            }
            失败返回 None
        """
        import json

        # 处理输入格式
        if isinstance(urls, str):
            url_list = [{"url": urls}]
        elif isinstance(urls, list):
            if len(urls) > 0 and isinstance(urls[0], str):
                url_list = [{"url": url} for url in urls]
            else:
                url_list = urls
        else:
            print("参数格式错误，需要传入字符串或列表")
            return None

        # 构建公共参数
        params = {
            'method': 'taobao.tbk.spread.get',
            'app_key': self.app_key,
            'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
            'format': 'json',
            'v': '2.0',
            'sign_method': 'md5',
            'partner_id': 'top-apitools',
        }

        # 添加业务参数 - requests 需要是 JSON 字符串
        params['requests'] = json.dumps(url_list, ensure_ascii=False)

        # 生成签名
        params['sign'] = self.generate_sign(params)

        # 发起请求
        try:
            response = requests.get(self.api_url, params=params, timeout=10)
            response.raise_for_status()

            result = response.json()

            # 检查是否有错误
            if 'error_response' in result:
                error = result['error_response']
                print(f"淘宝客 API 错误: {error.get('msg', '未知错误')}")
                print(f"错误代码: {error.get('code', 'N/A')}")
                print(f"子错误代码: {error.get('sub_code', 'N/A')}")
                return None

            # 返回成功结果
            if 'tbk_spread_get_response' in result:
                results = result['tbk_spread_get_response'].get('results', {}).get('tbk_spread', [])
                return results

            return None

        except requests.exceptions.RequestException as e:
            print(f"淘宝客 API 请求失败: {e}")
            return None
        except Exception as e:
            print(f"淘宝客 API 处理异常: {e}")
            return None


def create_taobao_api(app_key, app_secret):
    """创建淘宝客 API 实例"""
    return TaobaoAPI(app_key, app_secret)


def get_activity_info(app_key, app_secret, activity_material_id, adzone_id, **kwargs):
    """便捷函数：获取淘宝客活动信息"""
    api = TaobaoAPI(app_key, app_secret)
    return api.get_activity_info(activity_material_id, adzone_id, **kwargs)
