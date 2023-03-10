package get_available_slot

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/aaronangxz/AffiliateManager/orm"
	pb "github.com/aaronangxz/AffiliateManager/proto/affiliate"
	"github.com/aaronangxz/AffiliateManager/resp"
	"github.com/aaronangxz/AffiliateManager/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

type GetAvailableSlot struct {
	c         echo.Context
	ctx       context.Context
	requestId string
}

func New(c echo.Context) *GetAvailableSlot {
	d := new(GetAvailableSlot)
	d.c = c
	d.ctx = logger.NewCtx(d.c)
	logger.Info(d.ctx, "GetAvailableSlot Initialized")
	return d
}

func (g *GetAvailableSlot) GetAvailableSlotImpl() (*string, []*pb.BookingSlot, *resp.Error) {
	if err := g.verifyGetAvailableSlot(); err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	date := g.c.QueryParam("date")

	var err error
	dateTs, err := utils.GetDayTSWithDate(date)
	if err != nil {
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}
	dayStart, _ := utils.DayStartEndDate(time.Now().Unix())
	_, monthEnd := utils.MonthStartEndDate(time.Now().Unix())

	if dateTs < dayStart || dateTs > monthEnd {
		err := errors.New("time slot is out of range")
		logger.Error(g.ctx, err)
		return nil, nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_INVALID_PARAMS)
	}

	var slots []*pb.BookingSlot
	if err := orm.DbInstance(g.ctx).Raw(fmt.Sprintf("SELECT * FROM %v WHERE date = '%v'", orm.BOOKING_SLOTS_TABLE, date)).Scan(&slots).Error; err != nil {
		log.Error(err)
		return proto.String(date), nil, resp.BuildError(err, pb.GlobalErrorCode_ERROR_DATABASE)
	}
	return proto.String(date), slots, nil
}

func (g *GetAvailableSlot) verifyGetAvailableSlot() error {
	if g.c.QueryParam("date") == "" {
		return errors.New("invalid date")
	}
	return nil
}
