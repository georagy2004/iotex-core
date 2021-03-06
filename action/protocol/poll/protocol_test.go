// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package poll

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-election/test/mock/mock_committee"

	"github.com/iotexproject/iotex-core/action/protocol"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/test/mock/mock_chainmanager"
)

func TestNewProtocol(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sm := mock_chainmanager.NewMockStateManager(ctrl)
	committee := mock_committee.NewMockCommittee(ctrl)
	cfg := config.Default
	cfg.Consensus.Scheme = config.RollDPoSScheme
	p, err := NewProtocol(
		cfg,
		func(context.Context, string, uint64, time.Time, []byte) ([]byte, error) { return nil, nil },
		nil,
		committee,
		func(uint64) (time.Time, error) { return time.Now(), nil },
		sm,
	)
	require.NoError(err)
	require.NotNil(p)
}

func TestFindProtocol(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p, _, _, _, _ := initConstruct(ctrl)
	//if not registered
	re := protocol.NewRegistry()
	require.Nil(FindProtocol(re))

	//if registered
	require.NoError(p.Register(re))
	require.NotNil(FindProtocol(re))
}

func TestMustGetProtocol(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p, _, _, _, _ := initConstruct(ctrl)
	//if not registered
	re := protocol.NewRegistry()
	require.Panics(func() { MustGetProtocol(re) })

	//if registered
	require.NoError(p.Register(re))
	require.NotNil(FindProtocol(re))
}
