#!/usr/bin/env bash
set -euo pipefail

echo "checking forbidden L2, provider and x.go dependencies..."

MODULE_PATH="$(go list -m)"
DEPS="$(go list -deps ./...)"
FORBIDDEN_DEPS=(
  "github.com/bytechainx/x.go"
  "github.com/ZoneCNH/x.go"
  "github.com/bytechainx/foundationx"
  "github.com/ZoneCNH/foundationx"
  "github.com/bytechainx/redisx"
  "github.com/ZoneCNH/redisx"
  "github.com/bytechainx/kafkax"
  "github.com/ZoneCNH/kafkax"
  "github.com/bytechainx/postgresx"
  "github.com/ZoneCNH/postgresx"
  "github.com/bytechainx/taosx"
  "github.com/ZoneCNH/taosx"
  "github.com/bytechainx/ossx"
  "github.com/ZoneCNH/ossx"
  "github.com/bytechainx/market-engine"
  "github.com/ZoneCNH/market-engine"
  "github.com/redis/go-redis"
  "github.com/IBM/sarama"
  "github.com/segmentio/kafka-go"
  "github.com/jackc/pgx"
  "github.com/lib/pq"
  "gorm.io/driver/postgres"
  "github.com/taosdata/driver-go"
  "github.com/minio/minio-go"
  "github.com/aliyun/aliyun-oss-go-sdk"
)

for dep in "${FORBIDDEN_DEPS[@]}"; do
  if [[ "$dep" == "$MODULE_PATH" ]]; then
    continue
  fi
  if grep -Fq "$dep" <<<"$DEPS"; then
    echo "ERROR: testkitx must not depend on forbidden L2/provider dependency: $dep"
    exit 1
  fi
done

echo "checking forbidden production imports of testkitx..."

production_import_failed=0
while IFS= read -r -d '' file; do
  if grep -n --fixed-strings '"github.com/ZoneCNH/testkitx/pkg/testkitx' "$file"; then
    echo "ERROR: production file imports testkitx: $file"
    production_import_failed=1
  fi
done < <(find ./pkg ./internal -type f -name '*.go' ! -name '*_test.go' -print0)

if [[ "$production_import_failed" -ne 0 ]]; then
  exit 1
fi

echo "checking forbidden business terms..."

FORBIDDEN_TERMS=(
  "MacroRegime"
  "MarketRegime"
  "TradingSignal"
  "BTCUSDT"
  "ETHUSDT"
  "Kline"
  "OrderBook"
  "Position"
  "RiskGate"
)

for term in "${FORBIDDEN_TERMS[@]}"; do
  if grep -R --line-number --fixed-strings "$term" ./pkg ./internal --exclude-dir=.git; then
    echo "ERROR: forbidden business term found: $term"
    exit 1
  fi
done

echo "boundary check passed"
