<?php
/**
 * Author: 杨靖
 * Date: 2025/5/21
 * Time: 11:00
 */

namespace product\pay\service\jspay;

use common\helper\ArrayHelper;
use common\helper\FuncHelper;
use function DI\string;

/**
 * 江苏银行支付服务类
 */
class JsPayService
{
    protected $config = []; // 配置
    public $payCommonParams = []; // 公共请求报文

    public function __construct()
    {
        //配置
        $this->config = [
            'app_id' => "wxea66472a9f8b25af",        // appid
            'mch_id' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.mch_id'),        // 商户号
            'mch_name' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.mch_name'),        // 商户名称
            'public_key_path' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.public_key_path'),   // 江苏银行RSA公钥路径
            'pfx_path' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.pfx_path'),   // 赣邻通pfx文件路径
            'pfx_password' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.pfx_password'),   // 赣邻通pfx文件密钥
            'post_url' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.post_url'),   // 接口地址
            'master_account' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.master_account'),   // 主账号
            'notify_url' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.notify_url'),   // 异步通知地址
            'partnerId' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.partnerId'),   // 合作商ID
            'deviceNo' => ArrayHelper::getValue(\Yii::$app->params, 'jspay.wechat.deviceNo'),   // 终端编号
        ];

        // 公共请求报文
        $this->payCommonParams = [
            'createDate' => date('Ymd'), //报文发送日期
            'createTime' => date('His'), //报文发送时间
            'bizDate'    => date('Ymd'), //业务日期
            'msgID'      => md5(FuncHelper::uuid()), //报文标志号,唯一标识
            'svrCode'    => '', //服务代码,空即可
            'partnerId'  => $this->config['partnerId'], //合作商ID
            'channelNo'  => 'm', //发起渠道,固定为"m"
            'publicKeyCode' => '00', //公私钥对编号,固定为"00"
            'version'    => 'v1.0.0', //版本号,固定为"v1.0.0"
            'charset'    => 'utf-8', //字符集
            'signType'   => 'RSA', //签名类型
            'sign'       => '', //签名
            'deviceNo'   => $this->config['deviceNo'], //终端编号
        ];
    }

    /**
     * 统一支付入口
     * @param array $params 参数
     * @return array
     */
    public function pay($params = []): array
    {

        return match (ArrayHelper::getValue($params, 'type')) {
            'mini' => $this->wxMiniPay($params['order_sn'], $params['amount'], 'JSAPI',$params['open_id'],$params['app_id'] ?? $this->config['app_id']),
            'qrcode' => $this->wxQrcodePay($params['order_sn'], $params['amount']),
            default => [
                'ac' => 'err',
                'msg' => '不支持此支付方式',
            ],
        };
    }

    /**
     * 微信小程序支付、APP支付
     * @param string $order_sn 账单的唯一标识
     * @param float $amount 支付金额，单位：元
     * @param float $amount APP-APP 支付  JSAPI-小程序支付；
     * @param string $buyer_id 小程序支付时填openid；小程序支付时填openid；
     * @return array
     */
    public function wxMiniPay($order_sn, $amount, $tradeType, $openId, $app_id = ''): array
    {
        $data = [
            'service' => 'paymentWXPay',
            'totalFee' => $amount,  // 支付金额，单位：元
            'tradeType' => "JSAPI",  // JSAPI-小程序支付；
            'extfld2'  => $order_sn,  // 商户订单号
            'backUrl'  => $this->config['notify_url'], // 异步通知地址
            'mchIp'  => '0.0.0.0', // 订单生成的机器IP
            'extfld1' => $app_id, // 应用ID
            'openId' => $openId, // 小程序支付时填openid； 
        ];
    

        $data = array_merge($this->payCommonParams, $data);
        $data['sign'] = $this->sign($data);

        $result = $this->http($data);
        // APP支付、JSAPI支付参数保持一致
        if($tradeType == 'JSAPI'){
            $result['packAge'] = 'prepay_id='.$result['packAge'] ?? '';
        }

        if(empty($result['respCode'] ) || $result['respCode'] != '000000'){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行预下单失败，原因:'.($result['respMsg'] ?? '').';参数:'.json_encode($data,JSON_UNESCAPED_UNICODE));
            return [
                'ac' => 'err',
                'msg' => $result['respMsg'] ?? '',
            ];
        }

        return [
            'ac' => 'ok',
            'msg' => '',
            'data' => $result,
        ];
    }

    /**
     * 微信扫码支付
     * @param string $order_sn 账单的唯一标识
     * @param float $amount 支付金额，单位：元
     * @param string $mchIp 订单生成的机器IP
     * @return array
     */
    public function wxQrcodePay($order_sn, $amount, $mchIp = ''): array
    {
        $data = [
            'service' => 'dPay',
            'outTradeNo' => $order_sn, // 商户订单号
            'totalFee' => $amount, // 总金额 单位：元
            'qrValidTime' => 15, // 二维码有效时间，单位：分钟
            'backUrl'  => $this->config['notify_url'], // 异步通知地址
            'mchIp'  => $mchIp, // 订单生成的机器IP
            'deviceNo' => $this->config['deviceNo'], // 终端编号
        ];

        $data = array_merge($this->payCommonParams, $data);
        $data['sign'] = $this->sign($data);

        $result = $this->http($data);

        if($result['orderStatus'] != '1'){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行二维码生成失败，原因:'.($result['errMsg'] ?? '').';参数:'.json_encode($data,JSON_UNESCAPED_UNICODE));
            return [
                'ac' => 'err',
                'msg' => $result['errMsg'] ?? '',
            ];
        }

        return [
            'ac' => 'ok',
            'msg' => '',
            'pay_url' => $result['qrCode'], // 支付二维码链接
        ];
    }
   
  
    /**
     * 订单退款
     * @param string $order_sn 原支付商户订单号
     * @param string $refund_sn 退款单号
     * @param float $refund_amount 退款金额，单位：元
     * @return array
     */
    public function refund($order_sn, $refund_sn, $refund_amount): array
    {
        $data = [
            'service' => 'payRefund',
            'outTradeNo' => $order_sn, // 原支付商户订单号
            'outRefundNo' => $refund_sn, // 退款单号
            'refundAmt' => $refund_amount, // 退款金额 单位：元
            'deviceNo' => $this->config['deviceNo'], // 终端设备编号
        ];

        $data = array_merge($this->payCommonParams, $data);
        $data['sign'] = $this->sign($data);

        $result = $this->http($data);

        if($result['orderStatus'] != '1' && $result['orderStatus'] != '3'){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行订单退款失败，原因:'.($result['errMsg'] ?? '').';参数:'.json_encode($data,JSON_UNESCAPED_UNICODE));
            return [
                'ac' => 'err',
                'msg' => '江苏银行订单退款失败，原因:'.($result['errMsg'] ?? ''),
            ];
        }

        return [
            'ac' => 'ok',
            'msg' => '',
            'data' => $result,
        ];
    }

 
    /**
     * 订单查询
     * @param string $order_sn 商户订单号
     * @return array
     */
    public function query($order_sn): array
    {
        $data = [
            'service' => 'payCheck',
            'outTradeNo' => $order_sn, // 商户订单号
            'deviceNo' => $this->config['deviceNo']
        ];

        $data = array_merge($this->payCommonParams, $data);
        $data['sign'] = $this->sign($data);

        $result = $this->http($data);

        if($result['orderStatus'] != '1'){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行订单查询失败，原因:'.($result['errMsg'] ?? '').';参数:'.json_encode($data,JSON_UNESCAPED_UNICODE));
            return [
                'ac' => 'err',
                'msg' => '',
                'total_amount' => $result['totalFee'] ?? '0.00',
                'out_trade_no' => $result['outTradeNo'] ?? '',
                'data' => $result['orderStatus'],  // 1-成功、2-待查询、3-失败
            ];
        }

        return [
            'ac' => 'ok',
            'msg' => '',
            'total_amount' => $result['totalFee'] ?? '0.00',
            'out_trade_no' => $result['outTradeNo'] ?? '',
            'data' => $result['orderStatus'],  // 1-成功、2-待查询、3-失败
        ];
    }
 

    // 签名
    public function sign(array $params){
        $data = $params;
        unset($data['sign']);
        unset($data['signType']);

        $string = '';
        ksort($data);
        foreach ($data as $k => $v) {
            $string .= "&".$k."=".$v;
        }
        $string = substr($string, 1);

        $certs = array();
        openssl_pkcs12_read(file_get_contents(__DIR__.'/'.$this->config['pfx_path']), $certs, $this->config['pfx_password']);
        $privateKey = $certs['pkey'];

        // 使用私钥对数据进行签名
        if (!openssl_sign($string, $signature, $privateKey, OPENSSL_ALGO_SHA1)) {
            throw new \Exception('签名失败');
        }

        // 返回Base64编码后的签名
        return base64_encode($signature);
    }

    // 验签
    public function unSign(array $params){
        $signature = $params['sign'];
        $data = $params;
        unset($data['sign']);
        unset($data['signType']);

        $string = '';
        ksort($data);
        foreach ($data as $k => $v) {
            $string .= "&".$k."=".$v;
        }
        $string = substr($string, 1);

        $publicKey = file_get_contents(__DIR__.'/'.$this->config['public_key_path']); 
        // 验签
        $result = openssl_verify($string, base64_decode($signature), $publicKey, OPENSSL_ALGO_SHA1);
        if($result !== 1){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行异步回调验签失败:'.json_encode($params,JSON_UNESCAPED_UNICODE));
            throw new \Exception('验签失败');
        }

        return true;
    }

    // 发送请求
    private function http(array $params){
        try {
            ksort($params);
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行发送请求报文:'.json_encode($params,JSON_UNESCAPED_UNICODE));
            $opts = array(
                CURLOPT_TIMEOUT => 60,
                CURLOPT_RETURNTRANSFER => 1,
                CURLOPT_SSL_VERIFYPEER => false,
                CURLOPT_SSL_VERIFYHOST => false,
                CURLOPT_HTTPHEADER => ['Content-Type: text/html'],
                CURLOPT_URL => $this->config['post_url'],
                CURLOPT_POST => 1,
                CURLOPT_POSTFIELDS => urldecode(http_build_query($params)),
            );

            $ch = curl_init();
            curl_setopt_array($ch, $opts);

            $response = curl_exec($ch);
            $error = curl_error($ch);
            curl_close($ch);
            if ($error){
                throw new \Exception('请求失败,原因：'.$error);
            }
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行发送请求返回原始参数:'.$response);
            if(strpos($response, '&') === false){
                throw new \Exception('报错返回：'.$response);
            }

            // 分割字符串，获取参数数组
            $response = explode('&', $response);
            foreach ($response as $v){
                if(strpos($v,'sign=') !== false){
                    // 签名字段
                    $sign = substr($v, strlen('sign='), strlen($v));
                    $result['sign'] = $sign;
                    continue;
                }

                if(strpos($v,'qrCode=') !== false){
                    // 二维码字段
                    $qrCode = substr($v, strlen('qrCode='), strlen($v));
                    $result['qrCode'] = $qrCode;
                    continue;
                }

                if(strpos($v,'paySign=') !== false){
                    // 微信支付签名字段
                    $paySign = substr($v, strlen('paySign='), strlen($v));
                    $result['paySign'] = $paySign;
                    continue;
                }
                if(strpos($v,'datas=') !== false){
                    // PDF
                    $paySign = substr($v, strlen('datas='), strlen($v));
                    $result['datas'] = $paySign;
                    continue;
                }

                $map = explode("=", $v);
                $result[$map[0]] = $map[1];
            }
            if(empty($result)){
                throw new \Exception('请求失败');
            }

            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行发送请求返回参数:'.json_encode($result,JSON_UNESCAPED_UNICODE));
            if(empty($result['respCode'] ) || $result['respCode'] != '000000'){
                throw new \Exception($result['respMsg'] ?? '未知错误');
            }
        }catch (\Exception $e){
            \Yii::$app->prjLog->infoLog('jsPayError', '江苏银行发送请求失败，原因:'.$e->getMessage().';参数:'.json_encode($params,JSON_UNESCAPED_UNICODE));
            throw new \Exception('发送请求失败:'.$e->getMessage());
        }

        return $result;
    }
  

}