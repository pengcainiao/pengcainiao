package syncer

import (
	"context"
	"gitlab.com/a16624741591/zero/core/bloom"
	"os"

	"gitlab.com/a16624741591/zero/rest/httprouter"
)

var (
	filterMap = make(map[HoldingType]*bloom.Filter)
)

// BloomClient 布隆过滤器
func BloomClient(holdingType HoldingType) *bloom.Filter {
	var bitSize uint = 1024 * 1024 * 8 * 2 //默认2M
	//switch holdingType {
	//case HoldingTasks:
	//	bitSize = 1024 * 1024 * 8 * 5 //5M
	//case HoldingTaskDispatches:
	//	bitSize = 1024 * 1024 * 8 * 14 //14M
	//}
	// 单个键的大小为1M
	return bloom.New(Redis(), string(holdingType), bitSize)
}

func filterExists(filter *bloom.Filter, businessID string) httprouter.Response {
	if os.Getenv("BLOOM_FILTER") == "" {
		if b, err := filter.Exists(businessID); err != nil || !b {
			return httprouter.GetError(httprouter.ErrNotFoundCode, err)
		}
	}
	return httprouter.Success()
}

func IsInitialed(holdingType HoldingType) bool {
	return Redis().Exists(context.Background(), string(holdingType)).Val() == 1
}

//
//// TaskExists 判断事项是否存在
//func TaskExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingTasks]
//	if !ok {
//		filter = BloomClient(HoldingTasks)
//	}
//	return filterExists(filter, businessID)
//}
//
//// TaskDispatchExists 判断派发事项是否存在
//func TaskDispatchExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingTaskDispatches]
//	if !ok {
//		filter = BloomClient(HoldingTaskDispatches)
//	}
//	return filterExists(filter, businessID)
//}

// UserExists 判断用户是否存在
func UserExists(businessID string) httprouter.Response {
	filter, ok := filterMap[HoldingUsers]
	if !ok {
		filter = BloomClient(HoldingUsers)
	}
	return filterExists(filter, businessID)
}

//
//// FileExists 判断文件是否存在
//func FileExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingFiles]
//	if !ok {
//		filter = BloomClient(HoldingFiles)
//	}
//	return filterExists(filter, businessID)
//}
//
//// RecordExists 判断笔记是否存在
//func RecordExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingRecords]
//	if !ok {
//		filter = BloomClient(HoldingRecords)
//	}
//	return filterExists(filter, businessID)
//}
//
//// CommentExists 判断评论是否存在
//func CommentExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingComments]
//	if !ok {
//		filter = BloomClient(HoldingComments)
//	}
//	return filterExists(filter, businessID)
//}
//
//// ProjectExists 判断项目是否存在
//func ProjectExists(businessID string) httprouter.Response {
//	filter, ok := filterMap[HoldingProjects]
//	if !ok {
//		filter = BloomClient(HoldingProjects)
//	}
//	return filterExists(filter, businessID)
//}
