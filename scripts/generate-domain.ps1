param (
    [Parameter(Mandatory=$true, Position=0)]
    [string]$DomainName
)

$Pascal = $DomainName.Substring(0,1).ToUpper() + $DomainName.Substring(1)
$Lower = $DomainName.ToLower()
$Plural = $Lower + "s"

$rootDir = Split-Path -Parent $PSScriptRoot
if (!(Test-Path "$rootDir\internal")) {
    $rootDir = "C:\Go\GoProjectStructure"
}

$tplDir = "$PSScriptRoot\templates"
$utf8NoBom = New-Object System.Text.UTF8Encoding $false

Write-Host "⚡ Scaffolding clean production domain: $Pascal..." -ForegroundColor Cyan

function Process-Template($tplName, $targetPath) {
    $content = [System.IO.File]::ReadAllText("$tplDir\$tplName")
    $content = $content.Replace("__PASCAL__", $Pascal).Replace("__LOWER__", $Lower).Replace("__PLURAL__", $Plural)
    [System.IO.File]::WriteAllText($targetPath, $content, $utf8NoBom)
}

# 1. Model & DTOs
$modelFile = "$rootDir\internal\model\$Lower.go"
Process-Template "model.go.tpl" $modelFile
Write-Host "  [+] Created Model: $modelFile" -ForegroundColor Green

# 2. Repository (pgxpool queries)
$repoFile = "$rootDir\internal\repository\${Lower}_repository.go"
Process-Template "repository.go.tpl" $repoFile
Write-Host "  [+] Created Repository: $repoFile" -ForegroundColor Green

# 3. HTTP Handler (Chi endpoints & validation)
$handlerFile = "$rootDir\internal\handler\${Lower}_handler.go"
Process-Template "handler.go.tpl" $handlerFile
Write-Host "  [+] Created Handler: $handlerFile" -ForegroundColor Green

# 4. Database SQL Migrations
$timestamp = (Get-Date).ToString("yyyyMMddHHmmss")
$upMig = "$rootDir\internal\migration\${timestamp}_create_${Plural}_table.up.sql"
$downMig = "$rootDir\internal\migration\${timestamp}_create_${Plural}_table.down.sql"
Process-Template "migration_up.sql.tpl" $upMig
Process-Template "migration_down.sql.tpl" $downMig
Write-Host "  [+] Created Migrations: $upMig" -ForegroundColor Green

Write-Host ""
Write-Host "🚀 Domain $Pascal created successfully (Lean 3-Tier Production Architecture)!" -ForegroundColor Green
Write-Host "👉 Mount it in internal/handler/router.go in 2 steps:" -ForegroundColor Yellow
Write-Host "    ${Lower}Repo := repository.New${Pascal}Repository(db)" -ForegroundColor Cyan
Write-Host "    ${Lower}H := New${Pascal}Handler(${Lower}Repo)" -ForegroundColor Cyan
Write-Host ""
Write-Host "    r.Route(`/api/v1/$Plural`, func(r chi.Router) {" -ForegroundColor Cyan
Write-Host "        r.Get(`/`, ${Lower}H.List)" -ForegroundColor Cyan
Write-Host "        r.Get(`/{id}`, ${Lower}H.GetByID)" -ForegroundColor Cyan
Write-Host "        r.Post(`/`, ${Lower}H.Create)" -ForegroundColor Cyan
Write-Host "        r.Delete(`/{id}`, ${Lower}H.Delete)" -ForegroundColor Cyan
Write-Host "    })" -ForegroundColor Cyan
