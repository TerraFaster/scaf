# Create dist directory if it doesn't exist
$distPath = "dist"
if (-Not (Test-Path $distPath)) {
    New-Item -ItemType Directory -Path $distPath | Out-Null
}

# Build for multiple platforms
$targets = @(
    @{OS="windows"; ARCH="amd64"; Out="scaf-windows-amd64.exe"},
    @{OS="windows"; ARCH="arm64"; Out="scaf-windows-arm64.exe"},
    @{OS="linux";   ARCH="amd64"; Out="scaf-linux-amd64"},
    @{OS="linux";   ARCH="arm64"; Out="scaf-linux-arm64"},
    @{OS="darwin";  ARCH="amd64"; Out="scaf-darwin-amd64"},
    @{OS="darwin";  ARCH="arm64"; Out="scaf-darwin-arm64"}
)

foreach ($t in $targets) {
    Write-Host "Building for $($t.OS)/$($t.ARCH)..."
    $env:GOOS = $t.OS
    $env:GOARCH = $t.ARCH
    go build -o "$distPath\$($t.Out)" .
    Write-Host "Built $($t.Out)"
}

Write-Host "Build complete"