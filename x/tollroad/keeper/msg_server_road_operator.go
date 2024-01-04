package keeper

import (
	"context"
	"strconv"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateRoadOperator(goCtx context.Context, msg *types.MsgCreateRoadOperator) (*types.MsgCreateRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetRoadOperator(
		ctx,
		"0",
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	// Make sure that the new road operator has its ID taken from SystemInfo.
	// Have this ID returned by the message server function.
	// Make sure the next id in SystemInfo is incremented.
	// Emit an event with the expected type and attributes.

	systemInfo := &types.DefaultGenesis().SystemInfo
	saveInfo, isFound := k.GetSystemInfo(ctx)
	if isFound {
		systemInfo = &saveInfo
	}

	var roadOperator = types.RoadOperator{
		Creator: msg.Creator,
		Index:   strconv.Itoa(int(systemInfo.NextOperatorId)),
		Name:    msg.Name,
		Token:   msg.Token,
		Active:  msg.Active,
	}

	k.SetRoadOperator(
		ctx,
		roadOperator,
	)

	systemInfo.NextOperatorId = systemInfo.NextOperatorId + 1
	k.SetSystemInfo(ctx, *systemInfo)

	ctx.EventManager().EmitEvent(sdk.NewEvent("new-road-operator-created",
		sdk.NewAttribute("creator", roadOperator.Creator),
		sdk.NewAttribute("road-operator-index", roadOperator.Index),
		sdk.NewAttribute("name", roadOperator.Name),
		sdk.NewAttribute("token", roadOperator.Token),
		sdk.NewAttribute("active", strconv.FormatBool(msg.Active)),
	))

	return &types.MsgCreateRoadOperatorResponse{Index: roadOperator.Index}, nil
}

func (k msgServer) UpdateRoadOperator(goCtx context.Context, msg *types.MsgUpdateRoadOperator) (*types.MsgUpdateRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRoadOperator(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var roadOperator = types.RoadOperator{
		Creator: msg.Creator,
		Index:   msg.Index,
		Name:    msg.Name,
		Token:   msg.Token,
		Active:  msg.Active,
	}

	k.SetRoadOperator(ctx, roadOperator)

	return &types.MsgUpdateRoadOperatorResponse{}, nil
}

func (k msgServer) DeleteRoadOperator(goCtx context.Context, msg *types.MsgDeleteRoadOperator) (*types.MsgDeleteRoadOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetRoadOperator(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveRoadOperator(
		ctx,
		msg.Index,
	)

	return &types.MsgDeleteRoadOperatorResponse{}, nil
}
