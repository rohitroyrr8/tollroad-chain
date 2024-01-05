package keeper

import (
	"context"
	"fmt"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateUserVault(goCtx context.Context, msg *types.MsgCreateUserVault) (*types.MsgCreateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if isFound {
		return nil, sdkerrors.Wrap(types.ErrIndexSet, "index already set")
	}

	var userVault = types.UserVault{
		Owner:             msg.Creator,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	k.SetUserVault(
		ctx,
		userVault,
	)

	if userVault.Balance == 0 {
		return nil, types.ErrZeroTokens
	}

	creatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	err = k.bank.SendCoinsFromAccountToModule(ctx, creatorAddress, types.ModuleName, coinsOf(sdk.DefaultBondDenom, userVault.Balance))
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateUserVaultResponse{}, nil
}

func (k msgServer) UpdateUserVault(goCtx context.Context, msg *types.MsgUpdateUserVault) (*types.MsgUpdateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var userVault = types.UserVault{
		Owner:             msg.Creator,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	k.SetUserVault(ctx, userVault)

	creatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	if msg.Balance == 0 {
		return nil, types.ErrZeroTokens
	} else if msg.Balance > valFound.Balance {
		err := k.bank.SendCoinsFromAccountToModule(ctx, creatorAddress, types.ModuleName, coinsOf(sdk.DefaultBondDenom, msg.Balance-valFound.Balance))
		if err != nil {
			return nil, err
		}
	} else if msg.Balance < valFound.Balance {
		err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddress, coinsOf(sdk.DefaultBondDenom, valFound.Balance-msg.Balance))
		if err != nil {
			panic("bank error")
		}
	}

	return &types.MsgUpdateUserVaultResponse{}, nil
}

func (k msgServer) DeleteUserVault(goCtx context.Context, msg *types.MsgDeleteUserVault) (*types.MsgDeleteUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	fmt.Print("msg from msg_server_user_vault: ")
	fmt.Println(msg)
	valFound, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)

	creatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddress, coinsOf(sdk.DefaultBondDenom, valFound.Balance))
	if err != nil {
		panic("bank error")
	}
	fmt.Println("done")
	return &types.MsgDeleteUserVaultResponse{}, nil
}
