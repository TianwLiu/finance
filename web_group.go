package main

import "errors"


type WebGroup struct {
	GroupId       string   `json:"group_id"`
	MemberUserIds []string `json:"member_user_ids"`
}

var group map[string][]string

func init(){
	group= map[string][]string{
		"family": []string{"tianwei","guiqin","hongqi"},
	}

}

func addNewGroup(webGroup WebGroup)  error{
	if checkGroup(webGroup.GroupId){
		return errors.New("warning: new user trying to overwrite old group,execution cancel")
	}
	err:=updateWebGroup(webGroup)
	return err
}
func checkGroup(groupId string) bool{
	webGroup:=viewWebGroup(groupId)
	if webGroup.GroupId==groupId{
		return true
	}else{
		return false
	}
}
func addGroupMember(groupId string,newMemberId string) error{
	webGroup:=viewWebGroup(groupId)
	if webGroup.MemberUserIds!=nil{
		for _,memberId :=range webGroup.MemberUserIds{
			if memberId==newMemberId{
				return errors.New("the member has already been in the group")
			}
		}
		webGroup.MemberUserIds=append(webGroup.MemberUserIds, newMemberId)
	}else{
		webGroup=WebGroup{
			GroupId:       groupId,
			MemberUserIds: []string{groupId,newMemberId},
		}
	}
	err:=updateWebGroup(webGroup)
	return err
}
func delGroupMember(groupId string,oldMemberId string) error{
	webGroup:=viewWebGroup(groupId)

	if webGroup.MemberUserIds==nil{
		return errors.New("the group of given groupID not exist")
	}

	for k,memberId :=range webGroup.MemberUserIds{
		if memberId==oldMemberId{
			webGroup.MemberUserIds=append(webGroup.MemberUserIds[:k], webGroup.MemberUserIds[k+1:]...)
			err:=updateWebGroup(webGroup)
			return err
		}
	}
	return errors.New("the  member not exist in the group")

}
func getGroupMemberID(groupId string) ([]string,error){
	/*//fmt.Println(user[user_id])
	if group[groupId]==nil {
		return group[groupId],errors.New("can't find the users of specific group")
	}else{
		return group[groupId],nil
	}*/

	webGroup:=viewWebGroup(groupId)
	if webGroup.MemberUserIds==nil{
		return webGroup.MemberUserIds,errors.New("current group have no members")
	}else{
		return webGroup.MemberUserIds,nil
	}
}