Add-Type -AssemblyName System.Drawing

function Resize-Image($src, $dst, $w, $h) {
    $img = [System.Drawing.Image]::FromFile($src)
    $bmp = New-Object System.Drawing.Bitmap($w, $h)
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $g.DrawImage($img, 0, 0, $w, $h)
    $g.Dispose()
    $img.Dispose()
    $bmp.Save($dst, [System.Drawing.Imaging.ImageFormat]::Png)
    $bmp.Dispose()
}

$logoSrc = 'C:\Users\Admin\.gemini\antigravity\brain\7e736166-c580-4986-817e-1880dd3c8666\warp_speed_logo_1781188734340.png'
$logoDst = 'd:\Project\Personal\CLI tool\Screenshot\StoreLogo_300x300.png'
Resize-Image $logoSrc $logoDst 300 300

$posterSrc = 'C:\Users\Admin\.gemini\antigravity\brain\7e736166-c580-4986-817e-1880dd3c8666\warp_speed_poster_1781188748697.png'
$posterDst = 'd:\Project\Personal\CLI tool\Screenshot\StorePoster_720x1080.png'
Resize-Image $posterSrc $posterDst 720 1080
