package extract

import (
	"encoding/json"
	"fmt"
	internalCmd "github.com/aagu/go-i18n/pkg/cmd"
	"github.com/aagu/go-i18n/pkg/translation"
	"github.com/aagu/go-i18n/pkg/util"
	"github.com/spf13/cobra"
	lang "golang.org/x/text/language"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var sourceLang string
var outDir string
var mergeMode bool

func NewExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract <paths>",
		Short: "Extract all message definitions inside paths",
		Run:   extractFunc,
	}

	cmd.Flags().StringVar(&sourceLang, "source-lang", "en", "the language used by the extracted messages")
	cmd.Flags().StringVarP(&outDir, "out-dir", "o", ".", "the directory where messages files write to")
	cmd.Flags().BoolVarP(&mergeMode, "merge", "m", true, "merging existed message file inside out-dir")

	return cmd
}

func extractFunc(cmd *cobra.Command, args []string) {
	var paths []string
	if len(args) == 0 {
		paths = append(paths, ".")
	} else {
		paths = args
	}
	messages := make([]*translation.Message, 0)
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			// Ignore test files.
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			msgs, err := internalCmd.ExtractMessages(buf)
			if err != nil {
				return err
			}
			messages = append(messages, msgs...)
			return nil
		}); err != nil {
			cmd.PrintErrf("Extract message error: %v\n", err)
			os.Exit(1)
		}
	}

	if mergeMode {
		abs, err := util.GetAbsPath(outDir, false)
		if err == nil {
			parsedLang := lang.MustParse(sourceLang)

			path := filepath.Join(abs, fmt.Sprintf("%s.json", parsedLang.String()))
			readFile, err := ioutil.ReadFile(path)
			if err == nil {
				var existedMessages []*translation.Message
				if err := json.Unmarshal(readFile, &existedMessages); err == nil {
					messages = mergeMessages(messages, existedMessages)
				}
			}
		}
	}

	if err := writeFile(outDir, lang.MustParse(sourceLang), messages); err != nil {
		cmd.PrintErrf("Extract message error: %v\n", err)
		os.Exit(1)
	}
	cmd.Printf("Message extracted to %s\n", outDir)
}

func writeFile(outDir string, l lang.Tag, msgs []*translation.Message) error {
	abs, err := util.GetAbsPath(outDir, true)
	if err != nil {
		return err
	}
	path := filepath.Join(abs, fmt.Sprintf("%s.json", l.String()))
	content, err := json.Marshal(msgs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, content, 0644)
}

// mergeMessages will override dst with src
func mergeMessages(src, dst []*translation.Message) []*translation.Message {
	merged := make(map[string]*translation.Message)
	for idx := range dst {
		merged[dst[idx].ID] = dst[idx]
	}
	for idx := range src {
		merged[src[idx].ID] = src[idx]
	}
	var out []*translation.Message
	for _, message := range merged {
		out = append(out, message)
	}
	return out
}
