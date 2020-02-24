<#
.SYNOPSIS
    Script that builds the project for windows environemnt.    
#>
$OUTPUT_DIR = Join-Path $PSScriptRoot "output"
$DEFAULT_CONFIG_FILE = Join-Path $PSScriptRoot "app/conf/goaws.yaml"
$DEFAULT_CONFIG_DESTINATION = Join-Path $OUTPUT_DIR "conf"

Write-Host "Preparing to run build script..."

# Make sure output folder exists, #Otherwise clear the directory
if ((Test-Path $PSScriptRoot) -and !(Test-Path $OUTPUT_DIR)) {
    Write-Verbose -Message "Creating output directory..."
    New-Item -Path $OUTPUT_DIR -Type Directory | Out-Null
}
else {
    Write-Host "Clearing: $OUTPUT_DIR"
    Remove-Item -LiteralPath "$OUTPUT_DIR" -Force -Recurse
}

# Build GO project
Write-Host "Building the project..."
$buildCommand = "go build -o $OUTPUT_DIR/AbsorbLMS.GoAws.Service.exe app/cmd/goaws.go"
$invokeOutput = Invoke-Expression $buildCommand

if ($LASTEXITCODE -ne 0) {
    Write-Verbose -Message ($invokeOutput | Out-String)
    Throw "An error occurred while build the project."
}
Write-Host "Finished build"

# Copy default configuration
Write-Host "Copoying the default configuration file..."
if (!(Test-Path $DEFAULT_CONFIG_DESTINATION)) {
    New-Item -Path $DEFAULT_CONFIG_DESTINATION -Type Directory | Out-Null
}
Copy-Item -Path $DEFAULT_CONFIG_FILE -Recurse -Destination $DEFAULT_CONFIG_DESTINATION -Container

# Copy default configuration
Write-Host "Copoying the windows service installation script..."
Copy-Item -Path "InstallWindowsService.bat" -Destination $OUTPUT_DIR -Container

Write-Host "Package is ready in: $OUTPUT_DIR folder"