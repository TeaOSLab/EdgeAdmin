package users

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/userutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type FeaturesAction struct {
	actionutils.ParentAction
}

func (this *FeaturesAction) Init() {
	this.Nav("", "", "feature")
}

func (this *FeaturesAction) RunGet(params struct {
	UserId int64
}) {
	err := userutils.InitUser(this.Parent(), params.UserId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	featuresResp, err := this.RPC().UserRPC().FindAllUserFeatureDefinitions(this.AdminContext(), &pb.FindAllUserFeatureDefinitionsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	allFeatures := featuresResp.Features

	userFeaturesResp, err := this.RPC().UserRPC().FindUserFeatures(this.AdminContext(), &pb.FindUserFeaturesRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	userFeatureCodes := []string{}
	for _, userFeature := range userFeaturesResp.Features {
		userFeatureCodes = append(userFeatureCodes, userFeature.Code)
	}

	featureMaps := []maps.Map{}
	for _, feature := range allFeatures {
		featureMaps = append(featureMaps, maps.Map{
			"name":        feature.Name,
			"code":        feature.Code,
			"description": feature.Description,
			"isChecked":   lists.ContainsString(userFeatureCodes, feature.Code),
		})
	}
	this.Data["features"] = featureMaps

	this.Show()
}

func (this *FeaturesAction) RunPost(params struct {
	UserId int64
	Codes  []string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.User_LogUpdateUserFeatures, params.UserId)

	_, err := this.RPC().UserRPC().UpdateUserFeatures(this.AdminContext(), &pb.UpdateUserFeaturesRequest{
		UserId:       params.UserId,
		FeatureCodes: params.Codes,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
