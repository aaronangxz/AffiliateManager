package get_affiliate_details_by_id

import (
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/labstack/echo/v4"
)

type GetAffiliateDetailsById struct {
	c echo.Context
}

func New(c echo.Context) *GetAffiliateDetailsById {
	g := new(GetAffiliateDetailsById)
	g.c = c
	return g
}

func (g *GetAffiliateDetailsById) GetAffiliateDetailsByIdImpl() (*pb.AffiliateMeta, []*pb.ReferralDetails, *resp.Error) {
	if err := g.verifyGetAffiliateDetailsById(); err != nil {
		return nil, nil, err
	}

	id := g.c.Param("id")

	var (
		meta      *pb.AffiliateMeta
		refList   []*pb.ReferralDb
		finalList []*pb.ReferralDetails
	)

	//get affiliate meta
	if err := orm.DbInstance(g.c).Raw(orm.Sql2(), id).Scan(&meta).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	fmt.Println(meta)

	//get ref list
	if err := orm.DbInstance(g.c).Raw(orm.Sql3(), id).Scan(&refList).Error; err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	fmt.Println(refList)

	//TODO could be a bit redundant, dont need to show in this api call?
	for _, r := range refList {
		var bookingDetails *pb.BookingDetails
		if r.BookingId != nil {
			var b *pb.BookingDetailsDb
			if err := orm.DbInstance(g.c).Raw(orm.Sql4(), r.GetBookingId()).Scan(&b).Error; err != nil {
				return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
			}
			fmt.Println(b)
			bookingDetails = g.bookingDetailsDbToBookingDetails(b)
		}
		finalList = append(finalList, &pb.ReferralDetails{
			ReferralId:         r.ReferralId,
			AffiliateId:        r.AffiliateId,
			ReferralClickTime:  r.ReferralClickTime,
			ReferralStatus:     r.ReferralStatus,
			BookingDetails:     bookingDetails,
			ReferralCommission: r.ReferralCommission,
		})
	}
	fmt.Println(finalList)

	return meta, finalList, nil
}

func (g *GetAffiliateDetailsById) verifyGetAffiliateDetailsById() *resp.Error {
	if g.c.Param("id") == "" {
		return resp.BuildError(errors.New("invalid id"), pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	return nil
}

func (g *GetAffiliateDetailsById) bookingDetailsDbToBookingDetails(b *pb.BookingDetailsDb) *pb.BookingDetails {
	//TODO unmarshal CustomerInfo
	return &pb.BookingDetails{
		BookingId:        b.BookingId,
		BookingStatus:    b.BookingStatus,
		BookingDay:       b.BookingDay,
		BookingSlot:      b.BookingSlot,
		TransactionTime:  b.TransactionTime,
		PaymentStatus:    b.PaymentStatus,
		AdultTicketCount: b.AdultTicketCount,
		ChildTicketCount: b.ChildTicketCount,
		AdultTicketTotal: b.AdultTicketTotal,
		ChildTicketTotal: b.ChildTicketTotal,
		CustomerInfo:     nil,
	}
}
