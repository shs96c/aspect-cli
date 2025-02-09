/*
 * Copyright 2022 Aspect Build Systems, Inc.
 *
 * Licensed under the aspect.build Commercial License (the "License");
 * you may not use this file except in compliance with the License.
 * Full License text is in the LICENSE file included in the root of this repository.
 */

package test_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"aspect.build/cli/pkg/aspect/test"
	"aspect.build/cli/pkg/bazel/mock"
	"aspect.build/cli/pkg/ioutils"
	bep_mock "aspect.build/cli/pkg/plugin/system/bep/mock"
)

// Embrace the stutter :)
func TestTest(t *testing.T) {
	t.Run("test calls bazel test", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bzl := mock.NewMockBazel(ctrl)
		bzl.
			EXPECT().
			Spawn([]string{"test", "--bes_backend=grpc://127.0.0.1:12345"}).
			Return(0, nil)

		besBackend := bep_mock.NewMockBESBackend(ctrl)
		besBackend.
			EXPECT().
			Addr().
			Return("grpc://127.0.0.1:12345").
			Times(1)
		besBackend.
			EXPECT().
			Errors().
			Times(1)

		b := test.New(ioutils.Streams{}, bzl)
		g.Expect(b.Run([]string{}, besBackend)).Should(Succeed())
	})
}
