/**
 * 进度条
 *  $(".progress").progress({ selectorCompleteBackground: '#333' });
 */
;(function ($) {
    $.fn.progress = function (option) {
        //配置
        var setting = $.extend({
            selectorTotal: 'progress-bar', //总进度选择器类名
            selectorComplete: 'progress-complete', //已完成进度选择器类名
            selectorCompleteBackground: '#EFEFEF', //已完成进度选择器默认背景色
            selectorTotalLength: '250px', //总进度选择器默认宽度
            callback: function () {
            } //进度显示完后回调函数
        }, option);

        $(this).each(function (index, element) {
            //总数
            var totalNum = Number($(element).attr('data-total-num')) || 1;
            //总进度选择器宽度
            var totalLength = $(element).width() || setting.selectorTotalLength;

            //已完成进度选择器背景色
            var completeBackground = $(element).find('.' + setting.selectorComplete).css('background-color');

            //已完成数
            var completeNum = Number($(element).attr('data-complete-num')) || 1;
            //默认已完成进度选择器宽度
            var completeLength = totalLength;

            if (completeNum > totalNum)
                completeNum = totalNum;

            completeLength = completeNum / totalNum * totalLength;

            if(!completeBackground)
                $(element).find('.' + setting.selectorComplete).css('background-color', setting.selectorCompleteBackground).animate({width: Math.round(completeLength)}, "slow", setting.callback);
            else
                $(element).find('.' + setting.selectorComplete).animate({width: Math.round(completeLength)}, "slow", setting.callback);
        });

    }
})(jQuery);