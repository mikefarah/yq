$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
    PackageName    = 'yq'
    Url            = 'https://github.com/mikefarah/yq/releases/download/2.4.1/yq_windows_386.exe'
    Checksum       = '7ED72DF0565DC08481101DC5B5E00FC071791189541803F71BDD1BEB7FE90522'
    ChecksumType   = 'sha256'
    Url64bit       = 'https://github.com/mikefarah/yq/releases/download/2.4.1/yq_windows_amd64.exe'
    Checksum64     = 'BDFD2A00BAB3D8171EDF57AAF4E9A2F7D0395E7A36D42B07F0E35503C00292A3'
    ChecksumType64 = 'sha256'
    FileFullPath   = "$toolsDir\yq.exe"
}

Get-ChocolateyWebFile @packageArgs
