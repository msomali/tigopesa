// Package ussd
//
//How to use the NameCheckFunc
//This is an example struct that may attach business components like
//Database that has all the names for example
//type example struct {
//	Name string
//}
//
//func (e example) CustomNameCheckFuncByUser(ctx context.Context, request NameCheckRequest) (resp NameCheckResponse, err error)  {
//	fmt.Printf("request received is: %v\n",request)
//
//	resp = NameCheckResponse{
//		Result:    "customixe",
//		ErrorCode: "the error code",
//		ErrorDesc: "hey thetr",
//		Msisdn:    "jdjsudjd",
//		Flag:      "ndndhdhdhd",
//		Content:   "jdjdjdjd",
//	}
//
//	return resp, err
//}
//
//
//
//func test(handler NameCheckFunc)  {
//
//	req := NameCheckRequest{
//		Msisdn:              "071299lloouw",
//		CompanyName:         "startup",
//		CustomerReferenceID: "thereferenceid",
//	}
//
//	resp, err := handler.NameCheck(context.TODO(),req)
//
//	if err != nil {
//		fmt.Printf("%v\n",err)
//	}
//
//	fmt.Printf("%v\n",resp)
//}
//
//func main(){
//	e := example{Name: "exmaple struct reciver"}
//
//	test(e.CustomNameCheckFuncByUser)
//}
package ussd
