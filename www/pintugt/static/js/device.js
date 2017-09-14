var deviceWidth = document.documentElement.clientWidth;
if (deviceWidth > 750) deviceWidth = 750;
document.documentElement.style.fontSize = deviceWidth / 7.5 + 'px';
document.documentElement.style.margin = 0 + ' auto';
window.onload = window.onresize = function () {
    var deviceWidth = document.documentElement.clientWidth;
    if (deviceWidth > 750) deviceWidth = 750;
    document.documentElement.style.fontSize = deviceWidth / 7.5 + 'px';
    document.documentElement.style.margin = 0 + ' auto';
};
window.$$=function(els){
    var dom=document.querySelectorAll(els);
    if(dom.length===1){
        return dom[0];
    }else{
        return dom;
    }
}