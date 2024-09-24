<!--
 * @Author: FunctionSir
 * @Date: 2023-07-14 23:10:45
 * @LastEditTime: 2024-09-22 11:34:30
 * @LastEditors: FunctionSir
 * @Description: [A]nti [K]idnapping [B]eacon [P]roject
 * @FilePath: /AKBP/README.md
-->

# AKBP - [A]nti [K]idnapping [B]eacon [P]roject

Anti Kidnapping Beacon Project是受国际搜救卫星组织的信标启发的, 可用智能手机或其他可直接或间接连接到互联网的硬件作为信标的, 可自行架设服务器的, 反对绑架/虐待型"教育学校/矫正学校"的项目.  
愿每个人都可以不受绑架及虐待!  

## 正在进行重写

AKBP项目目前正在重写. 在正式宣布可用前, 仅可用于测试.  
**严禁**用于任何生产环境!  

## 务必注意

目前正在开发调试及测试阶段, 系统无法正常工作, 所有发来的数据包将可能不被认真对待.  

## 征求意见

1. 时间戳的单位问题.
2. 各(信标/服务器/HUB)ID的最大长度.

## 为什么有了这个项目?

众所周知, 某些所谓"网戒中心"/"教育学校"/"问题孩子学校"(勿对号入座)实在是恶心至极, 它们通过各种反人类的手段以求可以将一个所谓"问题少年"变成傀儡.  
这些孩子中有许多其实并不是什么问题少年, 而是只是有自己的思考, 没有完全听从家长的命令而已...  
这些被"教育"的未成年人往往深受其害. 而一旦他们被从家中带走, 比如被冒充成警察的所谓"教官"给带走, 未来会发生什么也就可以猜个大差不离了...  
如果是性少数人群, 且"家人"看不惯的话, 那么严酷的扭转"治疗"恐怕在所难免...  
找的对象不符合熊"家长"的心意被送进来的? 也不是没有.  
喜欢网络游戏? 害害! 只要"家长"想, 那也是会被带走的!  
总之, 这些机构干的事儿, 至少是把受害者从家里带到"学校"这个过程, 很难不让我联想到"Kidnapping", 也就是说绑架/诱拐/劫持/拐骗...  
很多时候一些有志之士想要救援, 然而, 他们可能会绝望的发现, 根本不知道人在哪儿...  
受到国际搜救卫星组织的启发, 鄙人觉得, 搞这么个信标项目很重要. 信息越多, 则诸如救援和取证, 都可能会更容易些.  

## 我想给我家孩子搞一个, 防范一般意义上的绑架, 可以么?

理论上来说, 可以. 毕竟, 呃... 就是完全可以.  

## 我... 需要购买特殊硬件?

不, 你完全不需要. 国际搜救卫星组织也许需要你拥有一个专用的信标, 而这个项目不同.  
很多人没有钱买一个自己的信标, 也有很多人如果即使买了一个信标也到不了手里(父母可能会拦截其快件)...
所以, 这个项目的目标之一, 就是构建一套可以用智能手机当作信标的系统. 当然, 如果你想的话, 电脑也不是不可以, 但是一些信息恐怕就不能自动获取了.  

## 我可以构建属于我自己的硬件信标么?

你当然可以. 这套系统是基于HTTP(S)的, 目前可以用HTTP请求来发送相关信息. 也就是说, 不管你是树莓派, 还是ESP32, 还是其他东西, 只要能进行特定请求, 就可以将信息送过来, 也就是说, 可以作为信标.  

## 架设一个服务器的需求

你需要AKBP-Server与Dokuwiki.  
建议您使用Linux系统, 因为目前本项目的所有程序均是在Linux(具体到发行版的话是Arch Linux)系统中开发与测试的.  
当然, Windows也是可以的, 但目前未经过测试.  
