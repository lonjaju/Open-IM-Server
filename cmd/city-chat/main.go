package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"Open_IM/internal/demo/register"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/utils"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.DebugMode)
	f, _ := os.Create("../logs/api-dev.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.Use(utils.CorsHandler())
	if config.Config.Prometheus.Enable {
		r.GET("/metrics", promePkg.PrometheusHandler())
	}
	authRouterGroup := r.Group("/demo|auth")
	{
		// 发送验证码
		authRouterGroup.POST("/code", register.SendVerificationCode)
		// 校验验证码
		authRouterGroup.POST("/verify", register.Verify)
		// 设置密码
		authRouterGroup.POST("/password", register.SetPassword)
		// 登录
		authRouterGroup.POST("/login", register.Login)
		// 重置密码
		authRouterGroup.POST("/reset_password", register.ResetPassword)
		// 检测user是否允许登录
		authRouterGroup.POST("/check_login", register.CheckLoginLimit)

		// 获取用户手机号
		authRouterGroup.POST("/one_click/get_phone", register.OneClickGetPhone)
	}

	{
		// 首页推荐
		// 进入首页提醒用户开启定位权限，否则不能使用，将显示错误缺省页。
		// 如果开启了定位，则通过位置信息、推荐方式（平台推荐或附近的位置）查询该地区的用户列表，如果该地区没有用户，则显示缺省页面
		// 1，查询条件：地址位置、推荐方式（1，平台推荐；2，根据地理位置距离）
		// 2，返回列表详情字段如下（根据用户会员时长排序，会员级别越高，开通时间越长，越靠前）
		// todo:

		// 圈子列表
		// 平台所有用户最新发布的动态消息列表
		// todo:

		// 我的消息
		// tab页面，分系统消息和用户聊天消息。消息状态为已读和未读，根据消息发布时间排序。消息列表第一行默认显示系统消息
		// todo:
	}

	{
		// 快捷消息管理
		// todo
	}

	{
		// todo:
		// 意见反馈
	}

	{
		// todo:
		// 公共接口
		// 1，用户基本信息
		// 2，关注用户
		// 3，上传文件/
	}

	//demoRouterGroup := r.Group("/auth")
	//{
	//	demoRouterGroup.POST("/code", register.SendVerificationCode)
	//	demoRouterGroup.POST("/verify", register.Verify)
	//	demoRouterGroup.POST("/password", register.SetPassword)
	//	demoRouterGroup.POST("/login", register.Login)
	//	demoRouterGroup.POST("/reset_password", register.ResetPassword)
	//	demoRouterGroup.POST("/check_login", register.CheckLoginLimit)
	//}

	//deprecated
	cmsRouterGroup := r.Group("/cms_admin")
	{
		cmsRouterGroup.POST("/generate_invitation_code", register.GenerateInvitationCode)
		cmsRouterGroup.POST("/query_invitation_code", register.QueryInvitationCode)
		cmsRouterGroup.POST("/get_invitation_codes", register.GetInvitationCodes)

		cmsRouterGroup.POST("/query_user_ip_limit_login", register.QueryUserIDLimitLogin)
		cmsRouterGroup.POST("/add_user_ip_limit_login", register.AddUserIPLimitLogin)
		cmsRouterGroup.POST("/remove_user_ip_limit_login", register.RemoveUserIPLimitLogin)

		cmsRouterGroup.POST("/query_ip_register", register.QueryIPRegister)
		cmsRouterGroup.POST("/add_ip_limit", register.AddIPLimit)
		cmsRouterGroup.POST("/remove_ip_Limit", register.RemoveIPLimit)
	}
	defaultPorts := config.Config.Demo.Port
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10004 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start demo api server address: ", address, ", OpenIM version: ", constant.CurrentVersion, "\n")
	go register.OnboardingProcessRoutine()
	go register.ImportFriendRoutine()
	err := r.Run(address)
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
	}
}
