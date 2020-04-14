package cli

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/kong"
	"github.com/willabides/bindown/v3"
	"github.com/willabides/bindown/v3/internal/configfile"
)

var kongVars = kong.Vars{
	"configfile_help":                 `file with bindown config`,
	"configfile_default":              `bindown.yml`,
	"cache_help":                      `directory downloads will be cached`,
	"install_help":                    `install a dependency`,
	"system_default":                  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	"system_help":                     `target system in the format of <os>/<architecture>`,
	"systems_help":                    `target systems in the format of <os>/<architecture>`,
	"checksums_help":                  `add checksums to the config file`,
	"config_format_help":              `formats the config file`,
	"config_validate_bin_help":        `name of the binary to validate`,
	"config_validate_help":            `validate that installs work`,
	"config_install_completions_help": `install shell completions`,
	"config_extract_path_help":        `output path to directory where the downloaded archive is extracted`,
	"install_force_help":              `force install even if it already exists`,
	"install_target_file_help":        `file to install`,
	"download_force_help":             `force download even if the file already exists`,
	"download_target_file_help":       `filename and path for the downloaded file. Default downloads to cache.`,
	"download_dependency_help":        `name of the dependency to download`,
	"download_help":                   `download a dependency but don't extract or install it`,
	"extract_dependency_help":         `name of the dependency to extract`,
	"extract_help":                    `download and extract a dependency but don't install it`,
	"extract_target_dir_help":         `path to extract to. Default extracts to cache.`,
}

var cli struct {
	Version            versionCmd                 `kong:"cmd"`
	Download           downloadCmd                `kong:"cmd,help=${download_help}"`
	Extract            extractCmd                 `kong:"cmd,help=${extract_help}"`
	Install            installCmd                 `kong:"cmd,help=${install_help}"`
	Format             fmtCmd                     `kong:"cmd,help=${config_format_help}"`
	AddChecksums       addChecksumsCmd            `kong:"cmd,help=${checksums_help}"`
	Validate           validateCmd                `kong:"cmd,help=${config_validate_help}"`
	ExtractPath        extractPathCmd             `kong:"cmd,help=${config_extract_path_help}"`
	InstallCompletions kong.InstallCompletionFlag `kong:"help=${config_install_completions_help}"`
	Configfile         string                     `kong:"type=path,help=${configfile_help},default=${configfile_default},env='BINDOWN_CONFIG_FILE'"`
	Cache              string                     `kong:"type=path,help=${cache_help},env='BINDOWN_CACHE'"`
	JSONConfig         bool                       `kong:"name=json,help='use json instead of yaml for the config file'"`
}

func configFile(ctx *kong.Context, filename string, noDefaultCache bool) *configfile.ConfigFile {
	config, err := configfile.LoadConfigFile(filename, noDefaultCache)
	ctx.FatalIfErrorf(err, "error loading config from %q", filename)
	if cli.Cache != "" {
		config.Cache = cli.Cache
	}
	return config
}

func newParser(kongOptions ...kong.Option) *kong.Kong {
	kongOptions = append(kongOptions,
		kong.Completers{
			"binpath": binPathCompleter,
			"bin":     binCompleter,
			"system":  systemCompleter,
		},
		kongVars,
		kong.UsageOnError(),
	)
	return kong.Must(&cli, kongOptions...)
}

//Run let's light this candle
func Run(args []string, kongOptions ...kong.Option) {
	parser := newParser(kongOptions...)

	kongCtx, err := parser.Parse(args)
	parser.FatalIfErrorf(err)
	err = kongCtx.Run()
	parser.FatalIfErrorf(err)
}

func init() {
	kongVars["extract_path_target_help"] = `file you want the extract path for`
}

type extractPathCmd struct {
	TargetFile string             `kong:"arg,required=true,help=${extract_path_target_help},completer=binpath"`
	System     bindown.SystemInfo `kong:"name=system,default=${system_default},help=${system_help},completer=system"`
}

func (d extractPathCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	binary := path.Base(d.TargetFile)

	extractDir, err := config.ExtractPath(binary, d.System)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(ctx.Stdout, extractDir)
	return err
}

func init() {
	kongVars["checksums_dep_help"] = `name of the dependency to update`
}

type addChecksumsCmd struct {
	Dependency string               `kong:"required=true,arg,help=${checksums_dep_help},completer=bin"`
	Systems    []bindown.SystemInfo `kong:"name=system,default=${system_default},help=${systems_help},completer=system"`
}

func (d *addChecksumsCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	err := config.AddChecksums(&bindown.ConfigAddChecksumsOptions{
		Dependencies: []string{filepath.Base(d.Dependency)},
		Systems:      d.Systems,
	})
	if err != nil {
		return err
	}
	return config.Write(cli.JSONConfig)
}

type fmtCmd struct {
	JSON bool `kong:"help='output json instead of yaml'"`
}

func (c fmtCmd) Run(ctx *kong.Context) error {
	cli.Cache = ""
	config := configFile(ctx, cli.Configfile, true)
	if config != nil {
		return config.Write(cli.JSONConfig)
	}
	return nil
}

type validateCmd struct {
	Dependency string               `kong:"required=true,arg,help=${config_validate_bin_help},completer=bin"`
	Systems    []bindown.SystemInfo `kong:"name=system,default=${system_default},completer=system"`
}

func (d validateCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	return config.Validate([]string{d.Dependency}, d.Systems)
}

type installCmd struct {
	Force      bool               `kong:"help=${install_force_help}"`
	TargetFile string             `kong:"required=true,arg,help=${install_target_file_help},completer=binpath"`
	System     bindown.SystemInfo `kong:"name=system,default=${system_default},help=${system_help},completer=system"`
}

func (d *installCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	binary := path.Base(d.TargetFile)
	binDir := path.Dir(d.TargetFile)
	return config.Install(binary, d.System, &bindown.ConfigInstallOpts{
		TargetDir: binDir,
		Force:     d.Force,
	})
}

type downloadCmd struct {
	Force      bool               `kong:"help=${download_force_help}"`
	System     bindown.SystemInfo `kong:"name=system,default=${system_default},help=${system_help},completer=system"`
	Dependency string             `kong:"required=true,arg,help=${download_dependency_help},completer=bin"`
	TargetFile string             `kong:"name=output,help=${download_target_file_help}"`
}

func (d *downloadCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	pth, err := config.DownloadDependency(d.Dependency, d.System, &bindown.ConfigDownloadDependencyOpts{
		TargetFile: d.TargetFile,
		Force:      d.Force,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx.Stdout, "downloaded %s to %s\n", d.Dependency, pth)
	return nil
}

type extractCmd struct {
	System     bindown.SystemInfo `kong:"name=system,default=${system_default},help=${system_help},completer=system"`
	Dependency string             `kong:"required=true,arg,help=${extract_dependency_help},completer=bin"`
	TargetDir  string             `kong:"name=output,help=${extract_target_dir_help}"`
}

func (d *extractCmd) Run(ctx *kong.Context) error {
	config := configFile(ctx, cli.Configfile, false)
	pth, err := config.ExtractDependency(d.Dependency, d.System, &bindown.ConfigExtractDependencyOpts{
		TargetDirectory: d.TargetDir,
		Force:           false,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx.Stdout, "extracted %s to %s\n", d.Dependency, pth)
	return nil
}
