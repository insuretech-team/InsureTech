param(
    [string]$OutputDir = "backend/inscore/secrets",
    [int]$KeySize = 2048
)

$ErrorActionPreference = "Stop"

if ($KeySize -lt 2048) {
    throw "KeySize must be >= 2048."
}

New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null

$privatePath = Join-Path $OutputDir "jwt_rsa_private.pem"
$publicPath = Join-Path $OutputDir "jwt_rsa_public.pem"

$rsa = [System.Security.Cryptography.RSA]::Create($KeySize)

try {
    $privateDer = $rsa.ExportRSAPrivateKey()
    $publicDer = $rsa.ExportSubjectPublicKeyInfo()

    function ConvertTo-Pem {
        param(
            [string]$Header,
            [byte[]]$Bytes
        )
        $b64 = [Convert]::ToBase64String($Bytes)
        $lines = ($b64 -split "(.{1,64})" | Where-Object { $_ -ne "" })
        $body = ($lines -join "`n")
        return "-----BEGIN $Header-----`n$body`n-----END $Header-----`n"
    }

    $privatePem = ConvertTo-Pem -Header "RSA PRIVATE KEY" -Bytes $privateDer
    $publicPem = ConvertTo-Pem -Header "PUBLIC KEY" -Bytes $publicDer

    [System.IO.File]::WriteAllText((Resolve-Path $OutputDir).Path + "\jwt_rsa_private.pem", $privatePem)
    [System.IO.File]::WriteAllText((Resolve-Path $OutputDir).Path + "\jwt_rsa_public.pem", $publicPem)

    Write-Host "Generated:"
    Write-Host "  $privatePath"
    Write-Host "  $publicPath"
}
finally {
    $rsa.Dispose()
}
