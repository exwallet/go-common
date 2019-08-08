/*
 * @Author: kidd
 * @Date: 1/29/19 6:22 PM
 *
 */
package console

import (
	"github.com/exwallet/go-common/gologger"
	"github.com/exwallet/go-common/goutil/gostring"
)

//PasswordPrompt 提示输入密码
//@param 是否二次确认
func InputPassword(prompt string, isConfirm bool, minLen int) (string, error) {
	var (
		confirm  string
		password string
		err      error
	)
	for {
		// 等待用户输入密码
		password, err = Stdin.PromptPassword(prompt)
		if err != nil {
			gologger.Error("unexpected error: %v\n", err)
			return "", err
		}
		if len(password) < minLen {
			gologger.Error("The length of the password is less than %d chars. Please re-enter it.\n", minLen)
			continue
		}
		// 二次确认密码
		if isConfirm {
			confirm, err = Stdin.PromptPassword("Confirm wallet password: ")
			if password != confirm {
				gologger.Error("The two password is not equal, please rre-enter it.\n")
				continue
			}
		}
		break
	}
	return password, nil
}

func InputString(prompt string, required bool) (string, error) {
	var (
		text string
		err  error
	)
	for {
		// 等待用户输入
		text, err = Stdin.PromptInput(prompt)
		if err != nil {
			gologger.Error("unexpected error: %v\n", err)
			return "", err
		}
		if len(text) == 0 && required {
			gologger.Error("Input can not be empty!\n")
			continue
		}
		break
	}
	return text, nil
}

func InputUint64(prompt string, isReal bool) (uint64, error) {
	var (
		num uint64
		ok  bool
	)
	for {
		// 等待用户输入参数
		line, err := Stdin.PromptInput(prompt)
		if err != nil {
			gologger.Error("unexpected error: %v\n", err)
			return 0, err
		}
		num, ok = gostring.NewString(line).UInt64()
		if !ok {
			gologger.Error("Input can not be empty, and must be greater than 0.\n")
			continue
		}
		if isReal && num <= 0 {
			gologger.Error("Input can not be empty, and must be greater than 0.\n")
			continue
		}
		break
	}
	return num, nil
}

func InputFloat64(prompt string, isReal bool) (float64, error) {
	var (
		num float64
		ok  bool
	)
	for {
		// 等待用户输入参数
		line, err := Stdin.PromptInput(prompt)
		if err != nil {
			gologger.Error("unexpected error: %v\n", err)
			return 0, err
		}
		num, ok = gostring.NewString(line).Float64()
		if !ok || (isReal && num <= 0) {
			gologger.Error("Input can not be empty, and must be greater than 0.\n")
			continue
		}
		break
	}
	return num, nil
}

func InputInt64(prompt string, isReal bool) (int64, error) {
	var (
		num int64
		ok  bool
	)
	for {
		// 等待用户输入参数
		line, err := Stdin.PromptInput(prompt)
		if err != nil {
			gologger.Error("unexpected error: %v\n", err)
			return 0, err
		}
		num, ok = gostring.NewString(line).Int64()
		if !ok || (isReal && num <= 0) {
			gologger.Error("Input can not be empty, and must be greater than or equal to 0.\n")
			continue
		}
		break
	}
	return num, nil
}
