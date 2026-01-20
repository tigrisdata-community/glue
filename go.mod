module github.com/tigrisdata-community/glue

go 1.25.5

tool (
	github.com/a-h/templ/cmd/templ
	github.com/caarlos0/pinata
	golang.org/x/tools/cmd/goimports
	golang.org/x/tools/cmd/stringer
	honnef.co/go/tools/cmd/staticcheck
)

require (
	github.com/aws/aws-sdk-go-v2 v1.41.0
	github.com/aws/aws-sdk-go-v2/config v1.32.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.95.0
	github.com/bwmarrin/discordgo v0.29.0
	github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
	github.com/gen2brain/avif v0.4.4
	github.com/gen2brain/webp v0.5.5
	github.com/go-faker/faker/v4 v4.7.0
	github.com/joho/godotenv v1.5.1
	github.com/openai/openai-go/v3 v3.16.0
	github.com/prometheus/client_golang v1.23.2
	github.com/tigrisdata/storage-go v0.2.0
)

require (
	charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251120230642-dcccabe2cd63 // indirect
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/a-h/parse v0.0.0-20250122154542-74294addb73e // indirect
	github.com/a-h/templ v0.3.977 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.5 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caarlos0/log v0.5.3 // indirect
	github.com/caarlos0/pinata v0.3.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charmbracelet/colorprofile v0.3.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20251120225753-26363bddd922 // indirect
	github.com/charmbracelet/x/ansi v0.11.1 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/cli/browser v1.3.0 // indirect
	github.com/clipperhouse/displaywidth v0.5.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/ebitengine/purego v0.8.3 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/telemetry v0.0.0-20260109210033-bd525da824e2 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)
