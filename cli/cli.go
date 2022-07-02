package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/pflag"
)

// Command : コマンド
type Command struct {
	// Name : コマンド名
	Name string
	// Alias : エイリアス
	Alias string
	// Short : 短いヘルプ文
	Short string
	// Long : 長いヘルプ文
	Long string
	// Example : サンプル
	Example string
	// Hidden : 表示しない
	Hidden bool

	// ValidateFunc : 引数のバリデーション
	ValidateFunc ValidateArgsFunc
	// SetFlagFunc : フラグをセット
	SetFlagFunc func(f *pflag.FlagSet)
	// RunFunc : 実行する関数
	RunFunc func(c *Command, f *pflag.FlagSet) error
	// HelpFunc : ヘルプ関数
	HelpFunc func(h string)

	// children : サブコマンド
	children map[string]*Command
}

// AddCommand : コマンドを追加
func (c *Command) AddCommand(newCommand *Command) {
	if c.children == nil {
		c.children = make(map[string]*Command)
	}

	c.children[newCommand.Name] = newCommand
}

// GetChildren : サブコマンドを取得
func (c *Command) GetChildren() []*Command {
	ls := []*Command{}

	for _, cmd := range c.children {
		if !cmd.Hidden {
			ls = append(ls, cmd)
		}
	}

	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Name < ls[j].Name
	})

	return ls
}

// GetChildrenNames : サブコマンド名の一覧を取得
func (c *Command) GetChildrenNames() []string {
	ls := []string{}

	for _, cmd := range c.GetChildren() {
		ls = append(ls, cmd.Name)
	}

	return ls
}

// newFlagSet : flagsetを生成
func (c *Command) newFlagSet() *pflag.FlagSet {
	f := pflag.NewFlagSet(c.Name, pflag.ContinueOnError)

	if c.SetFlagFunc != nil {
		c.SetFlagFunc(f)
	}

	f.BoolP("help", "h", false, fmt.Sprintf("help for %s", c.Name))

	return f
}

// find : サブコマンドを再帰的に検索
func find(cmd *Command, args []string) (*Command, []string) {
	for _, c := range cmd.GetChildren() {
		if args[0] != c.Name && args[0] != c.Alias {
			continue
		}

		if c.children == nil {
			return c, args[1:]
		}

		return find(c, args[1:])
	}

	return nil, nil
}

// getCommand : 引数から指定されたコマンドを取得
func getCommand(cmd *Command, args []string) (*Command, []string, error) {
	// コマンドを検索
	fCmd, fArgs := find(cmd, args)
	if fCmd == nil {
		return nil, nil, fmt.Errorf("command not found: %s", args[0])
	}

	// 引数チェック
	if fCmd.ValidateFunc != nil {
		if err := fCmd.ValidateFunc(fCmd, fArgs); err != nil {
			return nil, nil, err
		}
	}

	return fCmd, fArgs, nil
}

// Execute : 実行
func (c *Command) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("no argument")
	}

	cmd := c

	// 先頭がフラグでない
	if !strings.HasPrefix(args[0], "-") {
		var err error
		cmd, args, err = getCommand(cmd, args)
		if err != nil {
			return err
		}
	}

	// フラグを初期化
	f := cmd.newFlagSet()

	// パース
	if err := f.Parse(args); err != nil {
		return err
	}

	// ヘルプ
	if help, _ := f.GetBool("help"); help {
		cmd.help()
		return nil
	}

	return cmd.RunFunc(cmd, f)
}
