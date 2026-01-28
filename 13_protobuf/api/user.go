package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AVVKavvk/protobuf/proto"
	"github.com/labstack/echo/v4"
	proto_lib "google.golang.org/protobuf/proto"
)

var (
	userMap = make(map[int32]*proto.User) // Temporary user map
)

func CreateUser(ctx echo.Context) error {
	// 1. Read the raw binary body

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.ErrBadRequest
	}

	user := &proto.User{}

	if err = proto_lib.Unmarshal(body, user); err != nil {
		fmt.Printf("Error: %v", err)
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid protobuf data %v", err))
	}

	userMap[user.Id] = user

	res, err := proto_lib.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, "application/protobuf", res)

}

func GetUserById(ctx echo.Context) error {

	userId := ctx.Param("userId")
	if userId == "" {
		return ctx.JSON(http.StatusBadRequest, "UserId is required")
	}
	userIdInt, err := strconv.ParseInt(userId, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid UserId")
	}

	user, ok := userMap[int32(userIdInt)]
	if !ok {
		return ctx.JSON(http.StatusNotFound, "User not found")
	}

	res, err := proto_lib.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, "application/protobuf", res)

}
