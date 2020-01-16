/**
 * 倒计时
 <div class="time" data-start-time="1578450301"></div>
 <div class="time" data-start-time="1578448861"></div>
 <div class="time" data-start-time="1578450745"></div>
 <script type="text/javascript">
 $(function(){
        $('.time').timedown();
    })
 </script>
 */
;(function ($) {
    $.fn.extend({
        /**
         * 倒计时结果字符串
         * @param option
         * @returns {timeLast}
         */
        timeLast: function (option) {
            //当前时间戳
            var nowTime = Math.round((new Date()).getTime() / 1000)
            var defaultSetting = {
                startTime: nowTime, //默认当前时间戳
                timeShowSelector: $(this),
                dateLeft: '',
                dateRight: '天',
                hourLeft: '',
                hourRight: '时',
                minuteLeft: '',
                minuteRight: '分',
                secondLeft: '',
                secondRight: '秒',
                strStop: "已经结束"
            };
            var setting = $.extend(defaultSetting, option);

            //结束时间
            var ts = setting.startTime - nowTime;

            var tempts = ts - 1;
            var result = "";
            var timerResult = function timer() {
                if (tempts > 0) {
                    //天数
                    var dateNum = parseInt(tempts / 86400);
                    //小时数
                    countTime = tempts - dateNum * 86400;
                    var hourNum = parseInt(countTime / 3600);
                    //分数
                    countTime = countTime - hourNum * 3600;
                    var minuteNum = parseInt(countTime / 60);
                    //秒数
                    countTime = countTime - minuteNum * 60;
                    var secondNum = parseInt(countTime);
                    //减少一秒
                    tempts = tempts - 1;
                    //调整格式
                    var dateNumStr = (String(dateNum).length >= 2) ? dateNum : '0' + dateNum;
                    var hourNumStr = (String(hourNum).length >= 2) ? hourNum : '0' + hourNum;
                    var minuteNumStr = (String(minuteNum).length >= 2) ? minuteNum : '0' + minuteNum;
                    var secondNumStr = (String(secondNum).length >= 2) ? secondNum : '0' + secondNum;
                    result = "";
                    result += setting.dateLeft + dateNumStr + setting.dateRight;
                    result += setting.hourLeft + hourNumStr + setting.hourRight;
                    result += setting.minuteLeft + minuteNumStr + setting.minuteRight;
                    result += setting.secondLeft + secondNumStr + setting.secondRight;
                } else {
                    tempts = 0;
                    result = setting.strStop;
                    clearInterval(interval);
                }
                return result;
            }

            //显示结果
            function showResult() {
                var resultStr = timerResult;
                setting.timeShowSelector.html(resultStr);
            }

            var interval = setInterval(showResult, 1000);
            //支持链式操作
            return this;
        }
    });

    $.fn.extend({
        //倒计时
        timedown: function (option) {
            var setting = $.extend({}, option);
            $(this).each(function (index, element) {
                //开始时间，单位秒
                setting.startTime = $(element).attr('data-start-time') || 0;
                setting.timeShowSelector = $(element);
                $(element).timeLast(setting);
            })
        }
    });
})(jQuery);
