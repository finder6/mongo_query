package controllers

import (
	"container/list"
	"encoding/json"
	"fmt"
	"grxx_query/dao"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplName = "index.html"
}

func (this *MainController) Post() {
	param := this.GetString("p")
	if param != "" {
		number, err := strconv.ParseInt(param, 10, 64)
		var result []primitive.M
		params := list.New()
		if err != nil {
			result = doQuery(param, 0, params);
		} else {
			if(number < int64(2147483647)) {
				number32 := int32(number);
				result = doQuery(number32, 0, params);
			} else {
				result = doQuery(number, 0, params);
			}
		}
		// 过滤不要的字段
		for _, res := range result {
			delete(res, "_id")
			delete(res, "site")
			delete(res, "_")
		}
		var info string
		if len(result) != 0 {
			b, err := json.Marshal(result)
			if err != nil {
				fmt.Println("json.Marshal failed:", err)
				return
			}
			info = string(b)
		}
		if info == "" {
			info = "未找到数据"
		}
		
		this.Data["info"] = info
	}
	this.TplName = "index.html"
}

func doQuery(param interface{}, level int, params *list.List) []primitive.M {
	var filter primitive.M
	if level == 0 {
		filter = bson.M{"$or": []bson.M{bson.M{"qq":param},bson.M{"username":param},bson.M{"mobile":param},bson.M{"email":param}}}
	} else {
		filter = bson.M{"$or": []bson.M{bson.M{"qq":param},bson.M{"mobile":param},bson.M{"email":param}}}
	}
	params.PushFront(param)
	fmt.Println("list size = ", params.Len())
	result := dao.QueryMany(filter)
	if len(result) != 0 {
		if level < 10 {
			var resultAll []primitive.M
			for _, res := range result {
					if res["qq"] != nil {
						hasQuery := false
						for i := params.Front(); i != nil; i = i.Next() {
							if i.Value == res["qq"] {
								hasQuery = true
								break
							}
						}
						if hasQuery == false {
							resultAll = mergeResult(result, doQuery(res["qq"], level+1, params))
						}
					}
					if res["mobile"] != nil {
						hasQuery := false
						for i := params.Front(); i != nil; i = i.Next() {
							if i.Value == res["mobile"] {
								hasQuery = true
							}
						}
						if hasQuery == false {
							if resultAll == nil {
								resultAll = result
							}
							resultAll = mergeResult(resultAll, doQuery(res["mobile"], level+1, params))
						}
					}
					if res["email"] != nil && res["email"] != param {
						hasQuery := false
						for i := params.Front(); i != nil; i = i.Next() {
							if i.Value == res["email"] {
								hasQuery = true
							}
						}
						if hasQuery == false {
							if resultAll == nil {
								resultAll = result
							}
							resultAll = mergeResult(resultAll, doQuery(res["email"], level+1, params))
						}
					}
			}
			if resultAll == nil {
				resultAll = result
			}
			return resultAll
		}
	}
	return result
}

func mergeResult(a []primitive.M, b []primitive.M) []primitive.M{
	var d []primitive.M
	empty := true
	for _, k := range b {
		repeat := false
		for _, t := range a {
			if t["_id"] == k["_id"] {
				repeat = true
			}
		}
		if repeat == false {
			if empty {
				empty =  false
				d = append(a, k)
			} else {
				d = append(d, k)
			}
		}
	}
	if empty {
		return a;
	}
	return d
}