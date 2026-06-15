<?php
return [
    
    //江苏银行支付
    'jspay' => [
        //微信支付
        'wechat' => [
            'app_id' => 'wxea66472a9f8b25af',  //appid
            'master_account' => '10530188000186123', // 母账号
            'mch_id' => '20250211170842496001', // 商户号
            'mch_name' => '淮安高新商业管理有限公司', // 商户名称
            'bank_name' => '江苏银行', // 开户行
            'bank_no' => '313301099999', // 行号
            'public_key_path' => 'cert/jsbchinanew.cer',  //江苏银行RSA公钥路径
            'pfx_password' => '123456',  //
            'pfx_path' => 'cert/my_server.pfx',  // pfx文件路径
            'post_url' => 'https://mybank.jsbchina.cn:577/eis/merchant/merchantServices.htm',   // 接口地址
            'notify_url' => 'https://chaoshi.huaiangaoxin.top/pay-notify/jx-callback', //异步通知地址
            'partnerId' => '6dda59a1f57f409fb3b52adb8545847c', //合作商ID
            'deviceNo' => 'eb49d5f3adb9fd9f414726d70a60cfb9', // 终端编号
             
        ], 
    ],

    // 业主端小程序
    'xcx_yezhu' => [
        'app_id' => "wxea66472a9f8b25af",
        'secret' => '836c98c8453aa4456298aaae549287eb'
    ],
  
];