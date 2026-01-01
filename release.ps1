<#
.SYNOPSIS
    Automates the tagging and pushing process for SwiftStack.
.PARAMETER TagName
    The version tag (e.g., v1.0.0)
#>
param (
    [Parameter(Mandatory=$true)]
    [string]$TagName
)

# 1. Run Tests first
Write-Host "Running tests..." -ForegroundColor Cyan
#

# 2. Stage all changes
Write-Host "Staging changes..." -ForegroundColor Cyan
git add .

# 3. Commit
Write-Host "Committing..." -ForegroundColor Cyan
git commit -m "chore: release $TagName"

# 4. Create Tag
Write-Host "Creating tag $TagName..." -ForegroundColor Cyan
git tag -a $TagName -m "Release $TagName"

# 5. Push to GitHub (This triggers the GitHub Action)
Write-Host "Pushing to origin..." -ForegroundColor Cyan
git push origin main
git push origin $TagName

Write-Host "Done! Monitor your GitHub Actions for the build status." -ForegroundColor Green