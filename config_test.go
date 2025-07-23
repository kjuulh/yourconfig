package yourconfig_test

import (
	"testing"

	"github.com/kjuulh/yourconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("no types", func(t *testing.T) {
		type Config struct {
			SomeItem      string
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Zero(t, val)
	})

	t.Run("default tag, nothing set, no env set", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:""`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Zero(t, val)
	})

	t.Run("default tag (required=true), nothing set, no env set, err", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"required:true"`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.Error(t, err)
		require.Zero(t, val)

		assert.Equal(t, "config failed: field: SomeItem (env=SOME_ITEM) is not set and is required", err.Error())
	})

	t.Run("default tag (required=false), nothing set, no env set no error", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"required:false"`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Zero(t, val)
	})

	t.Run("env tag nothing set, no env set, no error", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"SOME_ITEM"`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Zero(t, val)
	})

	t.Run("default tag (required=true), nothing set, no env set, err", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"SOME_ITEM,required:true"`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.Error(t, err)
		require.Zero(t, val)

		assert.Equal(t, "config failed: field: SomeItem (env=SOME_ITEM) is not set and is required", err.Error())
	})

	t.Run("default tag (required), nothing set, no env set, err", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"SOME_ITEM,required"`
			someOtherItem string
			someBool      bool
		}

		val, err := yourconfig.Load[Config]()
		require.Error(t, err)
		require.Zero(t, val)

		assert.Equal(t, "config failed: field: SomeItem (env=SOME_ITEM) is not set and is required", err.Error())
	})

	t.Run("default tag private, trying to set, err", func(t *testing.T) {
		type Config struct {
			SomeItem      string
			someOtherItem string `cfg:"required:true"`
			someBool      bool
		}

		t.Setenv("SOME_OTHER_ITEM", "unsettable")

		val, err := yourconfig.Load[Config]()
		require.Error(t, err)
		require.Zero(t, val)

		assert.Equal(t, "config failed: field: someOtherItem is not settable", err.Error())
	})

	t.Run("env tag and env set, no error", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"required:true"`
			someOtherItem string
			someBool      bool
		}

		t.Setenv("SOME_ITEM", "some-item")

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Equal(t, "some-item", val.SomeItem)
	})

	t.Run("env tag (different name) and env set, no error", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"DIFFERENT_NAME,required:true"`
			someOtherItem string
			someBool      bool
		}

		t.Setenv("DIFFERENT_NAME", "some-item")

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Equal(t, "some-item", val.SomeItem)
	})

	t.Run("multiple env tag and env set, no error", func(t *testing.T) {
		type Config struct {
			SomeItem      string `cfg:"required:true"`
			SomeOtherItem string `cfg:"required:true"`
			someBool      bool
		}

		t.Setenv("SOME_ITEM", "some-item")
		t.Setenv("SOME_OTHER_ITEM", "some-other-item")

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Equal(t, "some-item", val.SomeItem)
		assert.Equal(t, "some-other-item", val.SomeOtherItem)
	})

	t.Run("required true, error returned", func(t *testing.T) {
		type Config struct {
			SomeItem      string
			SomeOtherItem string `cfg:"required:true"`
			someBool      bool
		}

		t.Setenv("SOME_OTHER_ITEM", "")

		val, err := yourconfig.Load[Config]()
		require.Error(t, err)
		require.Zero(t, val)

		assert.Equal(t, "config failed: field: SomeOtherItem (env=SOME_OTHER_ITEM) is not set and is required", err.Error())
	})

	t.Run("required false, no error returned", func(t *testing.T) {
		type Config struct {
			SomeItem      string
			SomeOtherItem string `cfg:"required:false"`
			someBool      bool
		}

		t.Setenv("SOME_OTHER_ITEM", "")

		val, err := yourconfig.Load[Config]()
		require.NoError(t, err)

		assert.Zero(t, val)
	})

}
