package administration

/*import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
	"tunnel/administration/model"
	"tunnel/config"
	"tunnel/common/logging"
)

func TestUserInfoService(t *testing.T) {
	logging.Initialize()
	admin, err := NewServerAdmin(config.Admin{
		Address: "",
		Port:    0,
		Https:   false,
		DBFile:  "E:\\Project\\tunnel\\tunn_server.db",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	newUserService(admin)
	service := userServiceInstance.infoService
	fmt.Println("---------------------------------------------------")
	fmt.Println("create")
	err = service.Create(&model.UserInfo{
		Id:         "0",
		Account:    "admin",
		Password:   "admin",
		Email:      "",
		Auth:       "",
		LastLogin:  0,
		LastLogout: 0,
		Created:    0,
		Updated:    0,
		FlowCount:  0,
		Disabled:   0,
		ConfigId:   "",
	})
	if err != nil {
		fmt.Println("error -> ", err)
		return
	} else {
		fmt.Println("success")
	}
	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("list1")
	list1, err := service.List()
	if err != nil {
		return
	}
	marshal1, err := json.Marshal(list1)
	if err != nil {
		return
	}
	fmt.Println("list --> ", string(marshal1))
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}

	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("create conflict")
	err = service.Create(&model.UserInfo{
		Id:       "0",
		Account:  "admin",
		Password: "admin",
	})
	if err != nil {
		fmt.Println("success -> ", err)
	} else {
		fmt.Println("failed")
		return
	}
	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("update")
	err = service.Save(&model.UserInfo{
		Id:         "0",
		Account:    "admin",
		Password:   "admin",
		Email:      "test",
		Auth:       "",
		LastLogin:  0,
		LastLogout: 0,
		Created:    0,
		Updated:    0,
		FlowCount:  0,
		Disabled:   0,
		ConfigId:   "",
	})
	if err != nil {
		fmt.Println("error -> ", err)
		return
	} else {
		fmt.Println("success")
	}
	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("list2")
	list2, err := service.List()
	if err != nil {
		return
	}
	marshal2, err := json.Marshal(list2)
	if err != nil {
		return
	}
	fmt.Println("list --> ", string(marshal2))
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}

	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("query")
	info := &model.UserInfo{
		Id:         "",
		Account:    "admin",
		Password:   "admin",
		Email:      "",
		Auth:       "",
		LastLogin:  0,
		LastLogout: 0,
		Created:    0,
		Updated:    0,
		FlowCount:  0,
		Disabled:   0,
		ConfigId:   "",
	}
	err = service.Query(info)
	if err != nil {
		fmt.Println("error -> ", err)
		return
	} else {
		m, _ := json.Marshal(info)
		fmt.Println("query --> ", string(m))

		if info.Email == "test" {
			fmt.Println("success")
		} else {
			fmt.Println("failed")
			return
		}
	}
	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("list3")
	list3, err := service.List()
	if err != nil {
		return
	}
	marshal3, err := json.Marshal(list3)
	if err != nil {
		return
	}
	fmt.Println("list --> ", string(marshal3))
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}

	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("delete")
	err = service.Delete(&model.UserInfo{
		Id:      "0",
		Account: "admin",
	})
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}
	info2 := &model.UserInfo{
		Id:         "0",
		Account:    "admin",
		Password:   "admin",
		Email:      "",
		Auth:       "",
		LastLogin:  0,
		LastLogout: 0,
		Created:    0,
		Updated:    0,
		FlowCount:  0,
		Disabled:   0,
		ConfigId:   "",
	}
	err = service.Query(info2)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("success")
		} else {
			fmt.Println("error -> ", err)
			return
		}
	}
	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("list4")
	list4, err := service.List()
	if err != nil {
		return
	}
	marshal4, err := json.Marshal(list4)
	if err != nil {
		return
	}
	fmt.Println("list --> ", string(marshal4))
	if err != nil {
		fmt.Println("error -> ", err)
		return
	}

	fmt.Println("---------------------------------------------------")

	fmt.Println("---------------------------------------------------")
	fmt.Println("query2")
	info3 := &model.UserInfo{
		Id:         "0",
		Account:    "admin",
		Password:   "admin",
		Email:      "",
		Auth:       "",
		LastLogin:  0,
		LastLogout: 0,
		Created:    0,
		Updated:    0,
		FlowCount:  0,
		Disabled:   0,
		ConfigId:   "",
	}
	err = service.Query(info3)
	if err != nil {
		fmt.Println("success -> ", err)
		return
	} else {
		m, _ := json.Marshal(info3)
		fmt.Println("query --> ", string(m))
	}
	fmt.Println("---------------------------------------------------")

}

//
// TestAdmin_Run
// @Description:
// @param t
//
func TestAdmin_Run(t *testing.T) {
	logging.Initialize()
	config.Current.Security = config.Security{
		CertPem: "../cmd/x509/cert.pem",
		KeyPem:  "../cmd/x509/key.pem",
	}
	admin, err := NewServerAdmin(config.Admin{
		Address: "",
		Port:    8080,
		Https:   true,
		DBFile:  "../tunn_server.db",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	admin.Run()
	time.Sleep(time.Hour * 24)
}*/
