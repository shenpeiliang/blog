#### 原理

1. web服务端生成一个唯一的key作为扫码登录事件ID，并设置过期时间
   
2. 根据这个key生成二维码信息，显示在网页登录页提供给APP扫码
   
3. 确认使用扫码登录方式后，前端需要每秒请求服务端询问这个key的状态，key过期则提示请刷新以获取有效的二维码
   
4. 用户使用APP扫码二维码（前提必须是登录状态），向服务端提供扫码获得的key和当前登录的uid
   
5. 前端请求获取到状态已经扫码，更新提示
   
6. APP请求查看这个key是否有效，有效则跳转到确认登录页，点击确认登录向服务端请求确认登录
   
7. 服务器根据uid信息登录保存会话返回APP提示已经授权成功，同时前端请求获取到状态已经确认登录，根据服务端返回的信息跳转到成功页面

#### 实现

1. 获取二维码
```
/**
* 前端登录二维码
*
*/
public function qr_show() {
    $this->load->model('sk_user_model');
    $this->load->library('sk_qrcode');

    //设置缓存信息：key、生成时间、状态、ip
    $key = $this->sk_user_model->qr_make_key(ip('int'));
    setcookie ( 'qr_code_key', $key, null, null, null, null, true );

    //这个跳转地址是提供给其他APP扫码跳转的地址，指定APP只获取key值
    $str = config_item('domain_login') . 'home/qr_scan?k=' . $key;

    header('Cache-Control: no-cache,no-store');
    die($this->sk_qrcode->png($str, false, 'L', 4, 1, '#F16522'));
}
```

2. 轮询请求二维码状态
```
/**
* 检查登录二维码
*/
public function qr_check() {
    if (empty($_COOKIE['qr_code_key'])) {
        ajax_response(true, 'COOKIE_NOT_SET');
    }

    $key = $_COOKIE['qr_code_key'];
    $this->load->model('sk_user_model');

    //key缓存是否存在，是否过期，状态是否合法，更新扫码状态
    $ret = $this->sk_user_model->qr_check($key, ip('int'));
    if ($ret) {
        // 登录成功，删除二维码cookie
        setcookie ( 'qr_code_key', false, null, null, null, null, true );
        ajax_response(true,'LOGIN_SUCCESS',$ret);
    }else{
        ajax_response(false,$this->sk_user_model->error);
    }
    
}
```

3. APP扫码登录

```
/**
* 二维码登录
*
*/
public function qr_scan() {
    $this->_filter_login();
    // 二维码key
    $key = $this->_param('key', null, true);
    // 是否确认
    $is_confirm = $this->_param('is_confirm', null, true);

    $this->load->model('sk_user_model');
    if ($is_confirm) {
        $ret = $this->sk_user_model->qr_confirm($key, $this->_user['id']);
    } else {
        $ret = $this->sk_user_model->qr_scan($key, $this->_user['id']);
    }
    $ret ? $this->_ok() : $this->_fail($this->sk_user_model->error);
}
```